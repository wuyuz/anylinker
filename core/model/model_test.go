package model

import (
	"testing"

	"anylinker/core/config"
	mylog "anylinker/core/utils/log"
)

func Test_countColums(t *testing.T) {
	querysql := `SELECT id,name,role,forbid,hashpassword FROM anyliner_user`
	anyliner := `SELECT count(*) FROM anyliner_user`
	gensql := gencountsql(querysql)
	if gensql != anyliner {
		t.Errorf("generate sql failed want getsql '%s',but gensql is '%s'", anyliner, gensql)
	}
}

func Test_ShowTable(t *testing.T) {
	config.Init("/Users/labulakalia/workerspace/golang/anyliner/core/config/core.toml")
	mylog.Init()
	InitDb()
	//InitRabc()
	// conn,err := db.GetConn(context.Background())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// rows,err := conn.QueryContext(context.Background(), "SELECT name FROM sqlite_master WHERE type ='table'")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// for rows.Next() {
	// 	var table string
	// 	rows.Scan(&table)
	// 	t.Log(table)
	// }
	//
	//enforcer := GetEnforcer()
	//pass, err := enforcer.Enforce("238397974042906624", "/api/v1/hostgroup", "POST")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(pass)
}
