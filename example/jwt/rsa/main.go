package main

import (
	"fmt"
	"github.com/zenfire-cn/webkit/jwt"
	"time"
)

type user struct {
	Name string
	Age  int
}

/**
 * @description: 可以先通过以下命令生成公钥和私钥
 *  openssl genrsa -out private.key 2048
 *  openssl rsa -in private.key -pubout > public.key
 */
func main() {
	// 初始化，设置token过期时间和公钥、私钥的文件路径
	jwt.RsaInit(time.Hour*2, "public.key", "private.key")

	// 传入对象（interface{}），生成token
	token, err := jwt.RsaGen(&user{Name: "Lorin", Age: 22})

	if err != nil {
		fmt.Println(err)
	}

	// 解析token
	res, err := jwt.RsaParse(token)
	fmt.Println(res, err)
}
