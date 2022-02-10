package user

import (
	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/utils/define"
	"anylinker/core/utils/helper"
	"anylinker/core/utils/resp"
	"context"
	"errors"
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
		return resp.JSON(c,resp.ErrBadRequest, nil)
	}

	hashpassword, err = utils.GenerateHashPass(ruser.Password) // 生成密码
	if err != nil {
		log.Error("GenerateHashPass failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}

	exist, err := model.Check(ctx, model.TBUser, model.Name, ruser.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if exist {
		return resp.JSON(c, resp.ErrUserNameExist, nil)
	}

	err = model.AddUser(ctx, ruser.Name, hashpassword,ruser.Role)
	if err != nil {
		log.Error("AddUser failed",zap.Error(err))
		return resp.JSON(c,resp.ErrInternalServer,nil)
	}
	return resp.JSON(c, resp.Success, nil)
}


// LoginUser login user
// @Summary login user
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/login [post]
// @Security BasicAuth
func LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	username, password, ok := helper.BasicAuth(c)
	if !ok {
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	token, err := model.LoginUser(ctx, username, password)
	if err != nil {
		log.Error("model.LoginUser failed", zap.Error(err))
		return c.SendStatus(resp.ErrUserPassword)
	}
	switch err := errors.Unwrap(err); err.(type) {
	case nil:
		resp.JSON(c, resp.Success, token)
	case define.ErrUserPass:
		resp.JSON(c, resp.ErrUserPassword, nil)
	case define.ErrForbid:
		resp.JSON(c, resp.ErrUserForbid, nil)
	default:
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	return nil
}

// GetUser Get User Info By Token
// @Summary get user info by token
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/info [get]
// @Security ApiKeyAuth
func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	uid := c.GetRespHeader("uid")
	// check uid exist
	exist, err := model.Check(ctx, model.TBUser, model.ID, uid)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		return resp.JSON(c, resp.ErrUserNotExist, nil)
	}

	user, err := model.GetUserByID(ctx, uid)
	if err != nil {
		log.Error("GetUserByID failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	user.Password = ""
	if user.Role == 2 {
		user.Roles = []string{"admin"}
	} else {
		user.Roles = []string{}
	}
	return resp.JSON(c, resp.Success, user)
}

// GetUsers get user info by token
// @Summary  get all users info
// @Tags User
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/all [get]
// @Security ApiKeyAuth
func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	var (
		q   define.Query
		err error
	)
	// TODO only admin
	err = c.QueryParser(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
		return c.SendStatus(resp.ErrInternalServer)
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	users, count, err := model.GetUsers(ctx, nil, q.Offset, q.Limit)
	if err != nil {
		log.Error("GetUsers failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	// remove password
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}
	return resp.JSON(c, resp.Success, users, count)
}

// ChangeUserInfo change user self config
// @Summary user change self's config info
// @Tags User
// @Description change self config,like email,wechat,dingphone,slack,telegram,password,remark
// @Produce json
// @Param User body define.ChangeUserSelf true "Change Self User Info"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/info [put]
// @Security ApiKeyAuth
func ChangeUserInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	newinfo := define.ChangeUserSelf{}
	err := c.QueryParser(&newinfo)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	if len(newinfo.Password) > 0 && len(newinfo.Password) < 8 {
		log.Error("password is short 8")
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	uid := c.Get("uid")
	if uid != newinfo.ID {
		log.Error("uid is error", zap.String("uid", uid), zap.String("infoid", newinfo.ID))
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	exist, err := model.Check(ctx, model.TBUser, model.UserName, newinfo.Name, newinfo.ID)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if exist {
		return resp.JSON(c, resp.ErrUserNameExist, nil)
	}
	err = model.ChangeUserInfo(ctx,
		uid,
		newinfo.Name,
		newinfo.Email,
		newinfo.WeChat,
		newinfo.DingPhone,
		newinfo.Telegram,
		newinfo.Password,
		newinfo.Remark)
	if err != nil {
		log.Error("ChangeUserInfo failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}

	return resp.JSON(c, resp.Success, nil)
}

// AdminChangeUser will change role,forbid,password,Remark
// @Summary admin change user info
// @Tags User
// @Description admin change user's role,forbid,password,remark
// @Produce json
// @Param User body define.AdminChangeUser true "Admin Change User"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/admin [put]
// @Security ApiKeyAuth
func AdminChangeUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	user := define.AdminChangeUser{}
	err := c.QueryParser(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	if len(user.Password) > 0 && len(user.Password) < 8 {
		log.Error("password is short 8")
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}
	// TODO only admin
	exist, err := model.Check(ctx, model.TBUser, model.ID, user.ID)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		return resp.JSON(c, resp.ErrUserNotExist, nil)
	}
	var role define.Role

	v := c.Get("role","")
	if v == "" {
		role = 3
	}
	if role != define.AdminUser {
		return resp.JSON(c, resp.ErrUnauthorized, nil)
	}

	err = model.AdminChangeUser(ctx, user.ID, user.Role, user.Forbid, user.Password, user.Remark)
	if err != nil {
		log.Error("AdminChangeUser failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}

	return resp.JSON(c, resp.Success, nil)
}

// AdminDeleteUser will delete user
// @Summary admin delete user
// @Tags User
// @Description admin delet user
// @Produce json
// @Param User body define.AdminChangeUser true "Admin Change User"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/admin [delete]
// @Security ApiKeyAuth
func AdminDeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	user := define.GetID{}
	err := c.QueryParser(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		return resp.JSON(c, resp.ErrBadRequest, nil)
	}

	var role define.Role
	v := c.Get("role","")
	if v == "" {
		role = 3
	}
	if role != define.AdminUser {
		return resp.JSON(c, resp.ErrUnauthorized, nil)
	}
	// 只能删除普通用户，不能删除admin用户
	userinfo, err := model.GetUserByID(ctx, user.ID)
	if err != nil {
		log.Error("GetUserByID failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if userinfo.Role == define.AdminUser {
		return resp.JSON(c, resp.ErrUnauthorized, nil)
	}
	// TODO only admin
	exist, err := model.Check(ctx, model.TBUser, model.ID, user.ID)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		return resp.JSON(c, resp.ErrUserNotExist, nil)
	}
	// 检查用户是否创建资源
	ok1, err := model.Check(ctx, model.TBTask, model.CreateByID, user.ID)
	if err != nil {
		log.Error("Check failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	ok2, err := model.Check(ctx, model.TBHostgroup, model.CreateByID, user.ID)
	if err != nil {
		log.Error("Check failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)

	}
	if ok1 || ok2 {
		return resp.JSON(c, resp.ErrDelUserUseByOther, nil)
	}
	err = model.DeleteUser(ctx, user.ID)
	if err != nil {
		log.Error("DeleteUser failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	return resp.JSON(c, resp.Success, nil)
}

// LogoutUser logout user
// @Summary logout user
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/logout [post]
// @Security BasicAuth
func LogoutUser(c *fiber.Ctx) error {
	return resp.JSON(c, resp.Success, nil)
}

// GetSelect return name,id
// @Summary return name,id
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/select [post]
// @Security BasicAuth
func GetSelect(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetNameID(ctx, model.TBUser)
	if err != nil {
		log.Error("model.GetNameID failed", zap.Error(err))
		return resp.JSON(c, resp.ErrInternalServer, nil)
	}
	return resp.JSON(c, resp.Success, data)
}