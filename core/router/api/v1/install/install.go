package install

import (
	"context"
	"github.com/gofiber/fiber/v2"

	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/define"
	"anylinker/core/utils/resp"
	//"anylinker/core/version"
	"go.uber.org/zap"
)

// QueryIsInstall query system is installed
func QueryIsInstall(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	isinstall, err := model.QueryIsInstall(ctx)
	if err != nil {
		log.Error("model.QueryIsInstall failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !isinstall {
		log.Debug("first running, need install...")
		return resp.JSON(c, resp.NeedInstall, nil)

	}
	return resp.JSON(c, resp.Success, nil)
}

// StartInstall install system
func StartInstall(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	isinstall, err := model.QueryIsInstall(ctx)
	if err != nil {
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if isinstall {
		return resp.JSON(c, resp.IsInstall, nil)
	}

	// get new create user
	adminuser := define.CreateAdminUser{}

	err = c.QueryParser(&adminuser)
	if err != nil {
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}

	err = model.StartInstall(ctx, adminuser.Name, adminuser.Password)
	if err != nil {
		log.Error("model.StartInstall", zap.Error(err))
		return resp.JSON(c, resp.ErrInstall, nil)
	}
	return resp.JSON(c, resp.Success, nil)
}

// QueryVersion query current version
func QueryVersion(c *fiber.Ctx) error {
	return resp.JSON(c, resp.Success, "V1.0.1")
}
