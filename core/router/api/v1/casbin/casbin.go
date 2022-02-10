package casbin

import (
	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/define"
	"anylinker/core/utils/resp"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetPermission get permission info by token
// @Summary  get all permission info
// @Tags permission
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/permission [get]
// @Security ApiKeyAuth
func GetPermission(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	var (
		q define.Query
		err error
	)

	err = c.QueryParser(&q)
	if  err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
		return c.SendStatus(resp.ErrInternalServer)
	}
	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}

	pers, count, err := model.GetPermissions(ctx, q.Offset,q.Limit )
	if err != nil {
		log.Error("GetPermission failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	return resp.JSON(c, resp.Success, pers, count)
}

// SetPermission get permission info by token
// @Summary  set new permission
// @Tags permission
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/set_permission [post]
// @Security ApiKeyAuth
func SetPermission(c *fiber.Ctx) error {
	ctx,cancel := context.WithTimeout(context.Background(),config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	per := define.SetPermission{}
	err := c.BodyParser(&per)
	if err != nil {
		log.Error("ParserBind failed", zap.Error(err))
		return resp.JSON(c,resp.ErrBadRequest, nil)
	}
	exist, err := model.Check(ctx,model.TBCasbin, model.All, per.P_type,per.V0,per.V1,per.V2)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if exist {
		return resp.JSON(c, resp.ErrUserNameExist, nil)
	}

	err = model.AddPermission(ctx,per.P_type,per.V0,per.V1,per.V2)
	if err != nil {
		log.Error("AddPermission failed",zap.Error(err))
		return resp.JSON(c,resp.ErrInternalServer,nil)
	}
	// 更新权限
	model.InitRabc()
	return resp.JSON(c, resp.Success, nil)
}

// AdminDeletePermission will delete user
// @Summary admin delete permission
// @Tags User
// @Description admin delet permission
// @Produce json
// @Param Permission body define.Setpermission true "Admin delete permission"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/set_permission [delete]
// @Security ApiKeyAuth
func DelPermission(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	per := define.Permission{}
	err := c.BodyParser(&per)
	if err != nil {
		log.Error("ParserBind failed", zap.Error(err))
		return resp.JSON(c,resp.ErrBadRequest, nil)
	}

	exist, err := model.Check(ctx,model.TBCasbin, model.All, per.P_type,per.V0.String,per.V1.String,per.V2.String)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if exist {
		err = model.DeletePermission(ctx,per.P_type,per.V0.String,per.V1.String,per.V2.String)
		if err != nil {
			log.Error("DeletePermission failed",zap.Error(err))
			return resp.JSON(c,resp.ErrInternalServer,nil)
		}
		model.InitRabc()
		return resp.JSON(c, resp.Success, nil)
	}
	return resp.JSON(c, resp.ErrUserNotExist, nil)
}