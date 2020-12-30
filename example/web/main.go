package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zenfire-cn/webkit/rest"
	"github.com/zenfire-cn/webkit/web"
	"go.uber.org/zap"
)

func main() {
	r := web.Init(gin.ReleaseMode, web.DefaultLog()) // 初始化：gin结合zap日志与lumberjack日志归档
	logger := zap.L()                                // 获取日志对象

	r.GET("/test/:num", func(c *gin.Context) {
		num := c.Param("num")

		logger.Info("params", zap.String("num", num))

		rest.Success(c, "Hello webkit, param is "+num)
	})

	r.Run(":8996")
}
