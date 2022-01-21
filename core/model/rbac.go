package model

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"anylinker/common/log"
	"anylinker/core/config"
	"go.uber.org/zap"
)

var (
	enforcer *casbin.Enforcer
)

// InitRabc init rabc
func InitRabc() {
	modeltext := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	dbcfg := config.CoreConf.Server.DB
	m, err := model.NewModelFromString(modeltext)
	if err != nil {
		log.Panic("NewModelFromString Err", zap.Error(err))
	}
	a, err := gormadapter.NewAdapter(dbcfg.Drivename, dbcfg.Dsn, true)
	if err != nil {
		log.Panic("NewAdapter Err", zap.Error(err))
	}

	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatal("InitRabc failed", zap.Error(err))
	}

}

// GetEnforcer get casbin auth
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
