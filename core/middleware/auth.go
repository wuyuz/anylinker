package middleware

import (
	"anylinker/common/jwt"
	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/resp"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const (
	tokenpre = "Bearer "
)

// CheckToken check token is valid
func CheckToken(token string) (string, string, bool) {
	claims, err := jwt.ParseToken(token)  // 校验token
	if err != nil || claims.UID == "" {
		log.Error("ParseToken failed", zap.Error(err))
		return "", "", false
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {  // 校验时间是否过期
		log.Error("Token is Expire", zap.String("token", token))
		return "", "", false
	}

	return claims.UID, claims.UserName, true
}

// 权限检查
func checkAuth(c *fiber.Ctx) (pass bool, err error) {
	token := strings.TrimPrefix(c.Get("Authorization"), tokenpre)

	if token == "" {
		err = errors.New("invalid token")
		return
	}
	uid, username, pass := CheckToken(token)  // 检验token
	if !pass {
		return false, errors.New("CheckToken failed")
	}
	c.Set("uid", uid)
	c.Set("username", username)

	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	ok, err := model.Check(ctx, model.TBUser, model.UID, uid)
	if err != nil {
		return false, err
	}

	if !ok {
		log.Error("Check UID not exist", zap.String("uid", uid))
		return false, nil
	}

	role, err := model.QueryUserRule(ctx, uid)
	if err != nil {
		log.Error("QueryUserRule failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	c.Set("role", role.String())

	requrl := c.Path()
	method := c.Method()
	enforcer := model.GetEnforcer()
	//fmt.Println(uid,requrl,method,enforcer)
	return enforcer.Enforce(role.String(), requrl, method)  // casbin鉴权
}

var excludepath = []string{"login", "logout", "install", "websocket","registry"}

// PermissionControl 权限控制middle
func PermissionControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			code = resp.Success
			err  error
		)
		if c.Path() == "/" {
			c.Next()
			return nil
		}
		// 白名单
		for _, url := range excludepath {
			if strings.Contains(c.Path(), url) {
				c.Next()
				return nil
			}
		}
		defer func() {
			c.Set("statuscode", strconv.Itoa(code))
		}()

		pass, err := checkAuth(c)  // 校验
		if err != nil {
			log.Error("checkAuth failed", zap.Error(err))
			code = resp.ErrUnauthorized
			goto ERR
		}
		if !pass {
			log.Error("checkAuth not pass ")
			code = resp.ErrUnauthorized
			goto ERR
		}

		c.Next()
		return nil

	ERR:
		// 解析失败返回错误
		c.Set("WWW-Authenticate", fmt.Sprintf("Bearer realm='%s'", resp.GetMsg(code)))
		resp.JSON(c, resp.ErrUnauthorized, nil)
		return nil
	}
}
