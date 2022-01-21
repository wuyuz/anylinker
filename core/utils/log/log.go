package log

import (
	"anylinker/common/log"
	"anylinker/core/config"
	"fmt"
	"os"
)

// Init init zap log
func Init() {
	logcfg := config.CoreConf.Log  // P诶之初始化

	err := log.InitLog(
		log.Path(logcfg.LogPath),
		log.Level(logcfg.LogLevel),
		log.Compress(logcfg.Compress),
		log.MaxSize(logcfg.MaxSize),
		log.MaxBackups(logcfg.MaxBackups),
		log.MaxAge(logcfg.MaxAge),
		log.Format(logcfg.Format),
	)
	if err != nil {
		fmt.Printf("InitLog failed: %v\n", err)
		os.Exit(1)
	}
}
