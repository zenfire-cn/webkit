package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zenfire-cn/webkit/jwt"
	"github.com/zenfire-cn/webkit/web"
	"time"
)

type user struct {
	Username string
	Password string
	Role     string
}

func main() {
	// 初始化 JWT，如果改用 Rsa 的初始化函数，登录时也使用 Rsa 的 token 生成函数
	jwt.Init(time.Hour*2, "mySecret")

	// 初始化 Web
	r := web.Init(gin.ReleaseMode, web.DefaultLog())

	// 挂载路由
	Router(r)

	r.Run(":8996")
}
