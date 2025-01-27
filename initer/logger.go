package initer

import (
	"log"

	"github.com/ztalab/ZACA/pkg/logger"
	"github.com/ztalab/ZACA/pkg/logger/redis_hook"
	"go.uber.org/zap/zapcore"

	"github.com/ztalab/ZACA/core"
)

func initLogger(config *core.Config) {
	conf := &logger.Conf{
		AppInfo: &logger.ConfigAppData{
			AppVersion: config.Version,
			Language:   "zh-cn",
		},
		Debug:  config.Debug,
		Caller: true,
	}
	if config.Debug {
		conf.Level = zapcore.DebugLevel
	} else {
		conf.Level = zapcore.InfoLevel
		conf.HookConfig = &redis_hook.HookConfig{
			Key:  config.Log.LogProxy.Key,
			Host: config.Log.LogProxy.Host,
			Port: config.Log.LogProxy.Port,
		}
	}
	if warn := logger.GlobalConfig(*conf); warn != nil {
		log.Print("[WARN] logger init error:", warn)
	}

	log.Print("[INIT] logger init success.")
}
