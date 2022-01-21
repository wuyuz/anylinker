package user

import (
	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/define"
	"anylinker/core/utils/resp"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/labulaka521/crocodile/common/utils"
	"go.uber.org/zap"
)

// RegistryUser new user
// @Summary registry new user
// @Tags User
// @Produce json
// @Param Registry body define.RegistryUser true "registry user"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/registry [post]
// @Security ApiKeyAuth
func RegistryUser(c *fiber.Ctx)  error {
	var (
		hashpassword  string
	)
	ctx,cancel := context.WithTimeout(context.Background(),config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	ruser := define.RegistryUser{}
	err := c.BodyParser(&ruser)
	if err != nil {
		log.Error("ParserBindUser failed", zap.Error(err))
		resp.JSON(c,resp.ErrBadRequest, nil)
		return nil
	}

	hashpassword, err = utils.GenerateHashPass(ruser.Password) // 生成密码
	if err != nil {
		log.Error("GenerateHashPass failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return nil
	}

	exist, err := model.Check(ctx, model.TBUser, model.Name, ruser.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return nil
	}
	if exist {
		resp.JSON(c, resp.ErrUserNameExist, nil)
		return nil
	}

	err = model.AddUser(ctx, ruser.Name, hashpassword,ruser.Role)
	if err != nil {
		log.Error("AddUser failed",zap.Error(err))
		resp.JSON(c,resp.ErrInternalServer,nil)
		return nil
	}
	resp.JSON(c, resp.Success, nil)
	return nil
}


