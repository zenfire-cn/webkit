package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zenfire-cn/commkit/logger"
	"go.uber.org/zap"
)

func Init(mode string, logOption *logger.Option) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.New()

	logger.Init(logOption)
	engine.Use(setGinLog(zap.L(), false))
	engine.Use(setGinRecovery(zap.L(), true))

	return engine
}

func DefaultLog() *logger.Option {
	return &logger.Option{
		Path:    "logs/webkit.log",
		Level:   "debug",
		MaxSize: 10,
		Json:    false,
		Std:     true, // 记录的同时是否也输出到控制台
	}
}
