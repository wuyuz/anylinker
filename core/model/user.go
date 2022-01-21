package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"anylinker/common/db"
	"anylinker/common/jwt"
	"anylinker/common/log"
	"anylinker/common/utils"
	"anylinker/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// LoginUser login user
func LoginUser(ctx context.Context, name string, password string) (string, error) {
	var (
		hashpassword string
		uid          string
		forbid       bool
	)
	loguser := `SELECT id,hashpassword,forbid FROM anyliner_user WHERE name=?`

	conn, err := db.GetConn(ctx)
	if err != nil {
		return "", fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, loguser)
	if err != nil {
		return "", fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, name).Scan(&uid, &hashpassword, &forbid)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("stmt.QueryRowContext Scan failed: %w", err)
	}
	if forbid {
		return "", fmt.Errorf(" failed: %w", define.ErrForbid{Name: name})
	}

	err = utils.CheckHashPass(hashpassword, password)
	if err != nil {
		return "", fmt.Errorf("utils.CheckHashPass failed: %w", define.ErrUserPass{Err: err})
	}
	token, err := jwt.GenerateToken(uid, name)
	if err != nil {
		return "", fmt.Errorf("jwt.GenerateToken failed: %w", err)
	}

	return token, nil
}

// AddUser add new user
func AddUser(ctx context.Context, name, hashpassword string, role define.Role) error {
	adduser := `INSERT INTO anyliner_user (
					id,
					name,
					hashpassword,
					role,
					forbid,
					createTime,
					updateTime
				)
				VALUES
				(?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, adduser)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	now := time.Now().Unix()
	id := utils.GetID()
	_, err = stmt.ExecContext(ctx, id, name, hashpassword, role, false, now, now)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}

	ok, err := enforcer.AddRoleForUser(id, role.String())
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("AddRoleForUser failed")
	}

	return nil
}

func getusers(ctx context.Context, uids []string, name string, offset, limit int) ([]define.User, int, error) {
	getsql := `SELECT
					id,
					name,
					role,
					forbid,
					hashpassword,
					email,
					wechat,
					dingphone,
					telegram,
					remark,
					createTime,
					updateTime
				FROM 
					anyliner_user`
	var (
		count int
		err   error
	)
	args := []interface{}{}
	users := []define.User{}
	if len(uids) > 0 {
		var querys = []string{}
		for _, uid := range uids {
			querys = append(querys, "id=?")
			args = append(args, uid)
		}
		getsql += " WHERE " + strings.Join(querys, " OR ")
	}
	if name != "" {
		getsql += " WHERE name=?"
		args = append(args, name)
	}

	if limit > 0 {
		count, err = countColums(ctx, getsql, args...)
		if err != nil {
			return users, 0, fmt.Errorf("countColums failed: %w", err)
		}

		getsql += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return users, 0, fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return users, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return users, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createTime int64
			updateTime int64
		)
		user := define.User{}
		err := rows.Scan(&user.ID,
			&user.Name,
			&user.Role,
			&user.Forbid,
			&user.Password,
			&user.Email,
			&user.WeChat,
			&user.DingPhone,
			&user.Telegram,
			&user.Remark,
			&createTime,
			&updateTime,
		)
		if err != nil {
			log.Error("Scan Err", zap.Error(err))
			continue
		}
		user.RoleStr = user.Role.String()
		user.CreateTime = utils.UnixToStr(createTime)
		user.UpdateTime = utils.UnixToStr(updateTime)
		if user.Role == define.AdminUser {
			user.Roles = []string{"admin"}
		} else {
			user.Roles = []string{}
		}
		users = append(users, user)
	}
	return users, count, nil
}

// GetUserByID get user by id
func GetUserByID(ctx context.Context, uid string) (*define.User, error) {
	userinfos, _, err := getusers(ctx, []string{uid}, "", 0, 0)
	if err != nil {
		return nil, fmt.Errorf("GerUser failed: %w", err)
	}
	if len(userinfos) != 1 {
		return nil, fmt.Errorf("Should get one user,but get total: %d", len(userinfos))
	}
	return &userinfos[0], nil
}

// GetUserByName get user by name
func GetUserByName(ctx context.Context, name string) (*define.User, error) {
	userinfos, _, err := getusers(ctx, nil, name, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("GerUser failed: %w", err)
	}
	if len(userinfos) != 1 {
		return nil, fmt.Errorf("Should get one user,but get total: %d", len(userinfos))
	}
	return &userinfos[0], nil
}

// GetUsers get all users info
func GetUsers(ctx context.Context, uids []string, offset, limit int) ([]define.User, int, error) {
	return getusers(ctx, uids, "", offset, limit)
}

// AdminChangeUser admin change user some column define.AdminChangeUser
// func AdminChangeUser(ctx context.Context, adminuser *define.AdminChangeUser) error {
func AdminChangeUser(ctx context.Context, id string, role define.Role, forbid bool, password, remark string) error {
	var (
		changeuser   string
		changerole   bool
		hashpassword string
	)
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	updateTime := time.Now().Unix()

	if password != "" {
		hashpassword, err = utils.GenerateHashPass(password)
		if err != nil {
			return fmt.Errorf("GenerateHashPass failed: %w", err)
		}
	} else {
		// get old user rolw
		userinfo, err := GetUserByID(ctx, id)
		if err != nil {
			return fmt.Errorf("GerUser failed: %w", err)
		}
		hashpassword = userinfo.Password
		if userinfo.Role != role {
			changerole = true
		}
	}

	// 普通管理员可以修改 password，role，forbid，
	changeuser = `UPDATE anyliner_user 
	SET role=?,
		forbid=?,
		hashpassword=?,
		updateTime=?,
		remark=?
	WHERE id=?`

	if changerole {
		// 修改权限表
		ok, err := enforcer.DeleteUser(id)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Delete user failed")
		}
		ok, err = enforcer.AddRoleForUser(id, role.String())
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("AddRoleForUser failed")
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			return fmt.Errorf("enforcer.LoadPolicy failed: %w", err)
		}
	}
	stmt, err := conn.PrepareContext(ctx, changeuser)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, role,
		forbid,
		hashpassword,
		updateTime,
		remark,
		id,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// ChangeUserInfo user change self's config define.ChangeUserSelf
// func ChangeUserInfo(ctx context.Context, id string, changeinfo *define.ChangeUserSelf) error {
func ChangeUserInfo(ctx context.Context, id, name, email, wechat, dingding, telegram, password, remark string) error {
	var (
		changeuser   string
		hashpassword string
	)
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	updateTime := time.Now().Unix()
	if password != "" {
		hashpassword, err = utils.GenerateHashPass(password)
		if err != nil {
			return fmt.Errorf("GenerateHashPass failed: %w", err)
		}
	} else {
		userinfo, err := GetUserByID(ctx, id)
		if err != nil {
			return fmt.Errorf("GerUser failed: %w", err)
		}
		hashpassword = userinfo.Password
	}
	changeuser = `UPDATE anyliner_user 
					SET hashpassword=?,
					    name=?,
						email=?,
						wechat=?,
						dingphone=?,
						telegram=?,
						updateTime=?,
						remark=? 
					WHERE
						id=?`

	stmt, err := conn.PrepareContext(ctx, changeuser)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, hashpassword,
		name,
		email,
		wechat,
		dingding,
		telegram,
		updateTime,
		remark,
		id,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// DeleteUser will delete user msg
func DeleteUser(ctx context.Context, id string) error {
	delsql := `DELETE FROM anyliner_user WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.Db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, delsql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}
