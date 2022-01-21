package router

import (
	"anylinker/core/middleware"
	"anylinker/core/router/api/v1/user"
	"anylinker/core/utils/define"
	"anylinker/core/config"
	"fmt"
	"errors"
	"go.uber.org/zap"
	"os"
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "anylinker/docs"

	//"google.golang.org/grpc"
	"net"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/soheilhy/cmux"
	"anylinker/common/log"
)

func NewHTTPRouter() *fiber.App {
	// 初始化router
	router := fiber.New()
	// router配置
	router.Use(pprof.New())
	router.Use(recover.New())

	router.Get("/docs/*", swagger.HandlerDefault)
	// 中间件
	router.Use(middleware.ZapLogger(),middleware.PermissionControl())

	v1 := router.Group("api/v1")
	ru := v1.Group("/user")
	{
		ru.Post("/registry", user.RegistryUser)


		ru.Post("/login", user.LoginUser)

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