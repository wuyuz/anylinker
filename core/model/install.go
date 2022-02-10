package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"anylinker/common/db"
	"anylinker/common/log"
	"anylinker/common/utils"
	"anylinker/core/config"
	"anylinker/core/utils/asset"
	"anylinker/core/utils/define"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"go.uber.org/zap"
)

var crcocodileTables = []string{
	TBHost,
	TBHostgroup,
	TBLog,
	TBNotify,
	TBOperate,
	TBTask,
	TBUser,
	TBCasbin,
}

// QueryIsInstall check table is create
func QueryIsInstall(ctx context.Context) (bool, error) {
	var querytable string
	needtables := []interface{}{}

	var queryname string
	params := []string{}
	drivename := config.CoreConf.Server.DB.Drivename
	if drivename == "sqlite3" {
		querytable = `SELECT count(*) FROM sqlite_master WHERE type="table" AND (`
		queryname = "name"
	} else if drivename == "mysql" {
		dbname := strings.Split(strings.Split(config.CoreConf.Server.DB.Dsn, "?")[0], "/")[1]
		needtables = append(needtables, dbname)
		querytable = `SELECT count(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=? AND (`
		queryname = "table_name"
	} else {
		return false, fmt.Errorf("unsupport drive type %s, only support sqlite3 or mysql", drivename)
	}

	for _, tbname := range crcocodileTables {
		needtables = append(needtables, tbname)
	}

	for i := 0; i < len(crcocodileTables); i++ {
		params = append(params, queryname+"=?")
	}
	querytable += strings.Join(params, " OR ")
	querytable += ")"
	var count int
	fmt.Println(querytable, params)
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	err = conn.QueryRowContext(ctx, querytable, needtables...).Scan(&count)
	if err != nil {
		log.Error("msg string", zap.Error(err))
		return false, fmt.Errorf("Scan failed: %w", err)
	}

	if count != len(crcocodileTables) {
		return false, nil
	}
	return true, nil
}

// StartInstall start install system
func StartInstall(ctx context.Context, username, password string) error {
	// create table
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}

	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}

	defer conn.Close()
	for _, tbname := range crcocodileTables {
		// crocodile_host
		var name string
		if tbname != TBCasbin {
			name = tbname[10:]
		} else {
			name = tbname
		}
		sqlfilename := "sql/" + name + ".sql"
		file, err := fs.Open(sqlfilename)
		if err != nil {
			log.Error("fs.Open failed", zap.String("filename", sqlfilename), zap.Error(err))
			continue
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error("ioutil.ReadAll failed", zap.Error(err))
			continue
		}
		var execsql string
		if config.CoreConf.Server.DB.Drivename == "sqlite3" {
			// sqlite3 TODO 的自增字段为AUTOINCREMENT
			execsql = strings.Replace(string(content), "AUTO_INCREMENT", "AUTOINCREMENT", -1)
			execsql = strings.Replace(string(content), "COMMENT", "--", -1)
		} else {
			execsql = string(content)
		}

		if tbname == TBCasbin {
			for _, sql := range strings.Split(execsql, ";\n") {
				if sql == "" {
					log.Warn("sql is empty string")
					continue
				}
				_, err = conn.ExecContext(context.Background(), sql)
				if err != nil {
					log.Error("conn.ExecContext failed", zap.Error(err), zap.String("sql", sql))
					return fmt.Errorf("conn.ExecContext failed: %w", err)
				}
			}
		} else {
			_, err = conn.ExecContext(ctx, execsql)
			if err != nil {
				log.Error("conn.ExecContext failed", zap.Error(err), zap.String("tbname", tbname))
				return fmt.Errorf("conn.ExecContext failed: %w", err)
			}
		}

		// wait second
		time.Sleep(time.Second / 2)
	}
	log.Debug("Success Run All crocodile Sql")

	// create admin user
	hashpassword, err := utils.GenerateHashPass(password)
	if err != nil {
		return fmt.Errorf("utils.GenerateHashPass failed: %w", err)
	}
	err = AddUser(ctx, username, hashpassword, define.AdminUser)
	if err != nil {
		return fmt.Errorf("AddUser failed: %w", err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Error("enforcer.LoadPolicy failed", zap.Error(err))
		return fmt.Errorf("enforcer.LoadPolicy failed: %w", err)
	}

	log.Debug("Success Install Crocodile")
	return nil
}
