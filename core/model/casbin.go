package model

import (
	"anylinker/common/db"
	"anylinker/core/utils/define"
	"context"
	"fmt"
	"go.uber.org/zap"
	"anylinker/common/log"

)

// GetHosts get all hosts
func GetPermissions(ctx context.Context, offset, limit int) ([]*define.Permission, int, error) {
	return getPermissionAll(ctx, "", "", offset, limit)
}


func getPermissionAll(ctx context.Context, role, t string, offset, limit int)([]*define.Permission, int, error){
	getsql := `SELECT
					p_type,
					v0,
					v1,
					v2,
					v3,
					v4,
					v5
				FROM
					casbin_rule`
	var(
		count int
		err error
	)
	arg := []interface{}{}
	pers := []*define.Permission{}

	// 角色查询
	if role != "" {
		getsql += " WHERE name=?"
		arg = append(arg, role)
	}
	if t != "" {
		getsql += " WHERE p_type=?"
		arg = append(arg, t)
	}

	if limit > 0 {
		count, err = countColums(ctx, getsql, arg...)
		if err != nil {
			return pers, 0, fmt.Errorf("countColums failed: %w", err)
		}

		getsql += " LIMIT ? OFFSET ?"
		arg = append(arg, limit, offset)
	}

	// 获取链接
	conn, err := db.GetConn(ctx)
	if err != nil {
		return pers, 0,fmt.Errorf("countColums failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return pers, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, arg...)
	if err != nil {
		return pers, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		p := define.Permission{}
		err := rows.Scan(
			&p.P_type,
			&p.V0,
			&p.V1,
			&p.V2,
			&p.V3,
			&p.V4,
			&p.V5,
			)
		if err != nil {
			log.Error("Scan Err", zap.Error(err))
			continue
		}
		pers=append(pers,&p)
	}
	return pers,count,nil
}