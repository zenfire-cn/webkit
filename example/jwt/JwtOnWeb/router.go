package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenfire-cn/webkit/jwt"
	"github.com/zenfire-cn/webkit/rest"
	"go.uber.org/zap"
)

/**
1. 生成 token 并返回
2. 在gin的路由中，使用 jwt.Auth 函数校验请求
3. 除 token 合法性以外，如果还需要校验其他信息（角色校验等），请实现 jwt.AuthCallback 的回调函数，并可以在 jwt.Auth 函数中按需传入参数
*/
func Router(r *gin.Engine) {
	// 1. 登录接口，生成 token 并返回
	r.POST("/login/:username/:password/:role", func(c *gin.Context) {
		username := c.Param("username") // 测试用例为了方便、直观，用户密码等信息直接获取url参数了
		password := c.Param("password")
		role := c.Param("role")

		if username == "admin" && password == "123456" {
			if token, err := jwt.Gen(&user{Username: username, Password: password, Role: role}); err == nil {
				rest.Success(c, token)
			} else {
				zap.L().Error(err.Error())
			}
		} else {
			rest.Error(c, "用户名或密码错误")
		}
	})

	// 2. 以 /api 分组的接口，加上了 jwt.Auth()，需要token校验通过才允许访问
	api := r.Group("/api", jwt.Auth())
	{
		api.GET("/hello", func(c *gin.Context) {
			jwtData, _ := c.Get("jwtData") // 如果要获取token解析出来的数据，c.Get("jwtData")
			rest.Success(c, jwtData)
		})
	}

	/**
	3. 如果需要校验其他信息（角色校验等），请实现 jwt.AuthCallback 中的回调函数
		jwtData map[string]interface{} -> token 解析出的数据
		validations ...interface{}     -> jwt.Auth 函数中传入的参数，例如：jwt.Auth("admin", "super")
	示例如下：
	*/
	jwt.AuthCallback(func(jwtData map[string]interface{}, validations ...interface{}) (bool, *rest.Rest) {
		fmt.Println(jwtData)
		if len(validations) > 0 {
			jwtRole := jwtData["Role"].(string) // 获取 token 中的角色信息
			needRole := validations[0].(string) // 获取参数中指定的角色
			if jwtRole == needRole {            // 角色符合，通过请求
				return true, nil
			} else { // 角色不符合，拦截请求并返回 rest 响应
				return false, &rest.Rest{Status: 500, Message: "您没有" + needRole + "权限以访问本接口"}
			}
		}
		return true, nil
	})

	// 实现了 jwt.AuthCallback 回调后，在jwt.Auth 中传入参数，就能在回调函数中按需使用
	auth := r.Group("/auth", jwt.Auth("super"))
	{
		auth.GET("/super", func(c *gin.Context) {
			jwtData, _ := c.Get("jwtData")
			rest.Success(c, jwtData)
		})
	}
}
