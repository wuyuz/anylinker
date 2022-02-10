package cmd

import (
	"anylinker/common/log"
	"anylinker/core/config"
	"anylinker/core/model"
	"anylinker/core/router"
	"anylinker/core/utils/define"
	mylog "anylinker/core/utils/log"

	//"anylinker/core/schedule"

	//"github.com/labulaka521/crocodile/core/router"
	//"github.com/labulaka521/crocodile/core/schedule"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"os"
)

// Server anylinker server
func Server() *cobra.Command {
	var (
		cfg string
	)

	cmdServer := &cobra.Command{
		Use:   "server",
		Short: "Start Run anylinker server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {  // 这里的的cfg会在后面调用时赋值，查看命令行数据
				cmd.Help()
				os.Exit(0)
			}
			// 传入配置文件地址，进行初始化，得到一个全局配置变量CoreConf
			config.Init(cfg)
			// 日志初始化
			mylog.Init()
			//alarm.InitAlarm()
			// 启动模型
			err := model.InitDb()
			if err != nil {
				log.Fatal("InitDb failed", zap.Error(err))
			}
			model.InitRabc()  // 初始化权限
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			lis, err := router.GetListen(define.Server)
			if err != nil {
				log.Fatal("listen failed", zap.Error(err))
			}
			// init alarm
			//err = schedule.Init2()
			if err != nil {
				log.Fatal("init schedule failed", zap.Error(err))
			}

			err = router.Run(define.Server, lis)
			if err != nil {
				log.Error("router.Run error", zap.Error(err))
			}
			return nil
		},
	}

	// 赋值变量
	cmdServer.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdServer
}

