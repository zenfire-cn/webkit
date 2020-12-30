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

func main() {
	// 初始化，设置token过期时间和自定义的秘钥
	jwt.Init(time.Hour*2, "mySecret")

	// 传入对象（interface{}），生成token
	token, err := jwt.Gen(&user{Name: "Lorin", Age: 22})

	if err != nil {
		fmt.Println(err)
	}

	// 解析token
	res, err := jwt.Parse(token)
	fmt.Println(res, err)
}
