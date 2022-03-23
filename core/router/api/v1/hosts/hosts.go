package hosts

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

// GetHost return all registry gost
// @Summary get all hosts
// @Tags Host
// @Description get all registry host
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host [get]
// @Security ApiKeyAuth
func GetHost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	var (
		q   define.Query
		err error
	)
	err = c.QueryParser(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	hosts, count, err := model.GetHosts(ctx, q.Offset, q.Limit)

	if err != nil {
		log.Error("GetHost failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}

	return resp.JSON(c, resp.Success, hosts, count)
}


// AddHost creat hosts
// @Summary creat hosts
// @Tags Host
// @Description creat hosts
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host [post]
// @Security ApiKeyAuth
func CreatHost(c *fiber.Ctx) error {
	ctx,cancel := context.WithTimeout(context.Background(),config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	host := define.AddHostRep{}
	err := c.BodyParser(&host)
	if err != nil {
		log.Error("ParserBind failed", zap.Error(err))
		return resp.JSON(c,resp.ErrBadRequest, nil)
	}
	exist, err := model.Check(ctx,model.TBHost, model.Addr, host.Addr)
	if err != nil {
	log.Error("IsExist failed", zap.Error(err))
	return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if exist {
	return resp.JSON(c, resp.ErrUserNameExist, nil)
	}

	//hashpassword, err = utils.GenerateHashPass(host.Password) // 生成密码
	//fmt.Println("xx",host,host.Addr)
	err = model.AddNewHost(ctx,host.Addr,host.HostName,host.UserName,host.Password)

	if err != nil {
		log.Error("AddUser failed",zap.Error(err))
		return resp.JSON(c,resp.ErrInternalServer,nil)
	}
	return resp.JSON(c, resp.Success, nil)
}