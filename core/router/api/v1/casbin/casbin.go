package casbin

import (
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/define"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"anylinker/common/log"
	"anylinker/core/utils/resp"

)

// GetPermission get permission info by token
// @Summary  get all permission info
// @Tags permission
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/permission/all [get]
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
