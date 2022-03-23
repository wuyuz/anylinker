package router

import (
	"anylinker/core/config"
	"anylinker/core/middleware"
	"anylinker/core/router/api/v1/casbin"
	"anylinker/core/router/api/v1/hosts"
	"anylinker/core/router/api/v1/install"
	"anylinker/core/router/api/v1/user"
	"anylinker/core/router/api/v1/ws"
	"anylinker/core/utils/define"
	_ "anylinker/docs"
	"errors"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
	"os"

	"anylinker/common/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/soheilhy/cmux"
	//"google.golang.org/grpc"
	"net"
)

func NewHTTPRouter() *fiber.App {
	// 初始化router
	router := fiber.New()
	// router配置
	router.Use(pprof.New())
	router.Use(recover.New())
	// 跨域
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: "http://localhost:9528",
			AllowHeaders:  "Origin, Content-Type, Accept,Authorization",
			AllowCredentials: true,
		},
	))

	// websocket
	//router.Use(func(c *fiber.Ctx) error {
	//	if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
	//		return c.Next()
	//	}
	//	return c.SendStatus(fiber.StatusUpgradeRequired)
	//})
	router.Get("/ws",websocket.New(ws.WsSsh))
	router.Get("/docs/*", swagger.HandlerDefault)

	// 中间件
	router.Use(middleware.PermissionControl(),middleware.ZapLogger())

	v1 := router.Group("dev/api/v1")
	ru := v1.Group("/user")
	{
		ru.Post("/registry", user.RegistryUser)
		ru.Post("/login",user.LoginUser)
		ru.Get("/info", user.GetUser)
		ru.Get("/all", user.GetUsers)             // only admin
		ru.Put("/admin", user.AdminChangeUser)    // only admin  // 管理员修改了某某用户
		ru.Delete("/admin", user.AdminDeleteUser) // only admin  // 管理员删除普通用户
		ru.Post("/logout", user.LogoutUser) // 某某注销登陆

		ru.Get("/permission",casbin.GetPermission)
		ru.Post("/set_permission",casbin.SetPermission)
		ru.Delete("/set_permission",casbin.DelPermission)
	}
	ri := v1.Group("/install")
	{
		ri.Get("/status", install.QueryIsInstall)
		ri.Post("", install.StartInstall)
		ri.Get("/version", install.QueryVersion)
	}
	rh := v1.Group("/hosts")
	{
		rh.Get("", hosts.GetHost)
		rh.Post("", hosts.CreatHost)
		//rh.Put("/stop", host.ChangeHostState)
		//rh.Delete("", host.DeleteHost)
		//rh.Get("/select", host.GetSelect)
	}

	return router
}

// GetListen get listen addr by server or client
func GetListen(mode define.RunMode) (net.Listener, error) {
	var (
		addr string
	)
	switch mode {
	case define.Server:
		if os.Getenv("PORT") != "" {
			addr = ":" + os.Getenv("PORT")
		} else {
			addr = fmt.Sprintf(":%d", config.CoreConf.Server.Port)
		}

	case define.Client:
		addr = fmt.Sprintf(":%d", config.CoreConf.Client.Port)

	default:
		return nil, errors.New("Unsupport mode")
	}
	lis, err := net.Listen("tcp", addr)
	return lis, err
}


func Run(mode define.RunMode, lis net.Listener) error {
	var (
		//gRPCServer *grpc.Server
		httpServer *fiber.App
		//err        error
		m          cmux.CMux
	)


	m = cmux.New(lis)
	if mode == define.Server {
		httpServer = NewHTTPRouter()
		httpL := m.Match(cmux.HTTP1Fast())
		go httpServer.Listener(httpL)
		log.Info("start run http server", zap.String("addr", lis.Addr().String()))
	}

	return m.Serve()
}