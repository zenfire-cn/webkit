package test

import (
	"fmt"
	"github.com/zenfire-cn/webkit/jwt"
	"testing"
	"time"
)

type testStruct struct {
	Name string
	Age  int
}

func TestGenAndParse(t *testing.T) {
	// 初始化
	jwt.Init(time.Hour*2, "mySecret")
	// 生成token
	token, err := jwt.Gen(&testStruct{"Lorin", 22})
	fmt.Println(token, err)
	// 解析token
	res, err := jwt.Parse(token)
	fmt.Println(res, err)
}

func TestData(t *testing.T) {
	data, err := jwt.Data("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7Ik5hbWUiOiJMb3JpbiIsIkFnZSI6MjJ9LCJleHAiOjE2MTE3MTgwNjV9.OI2u19Z5M6b5U1-vtI5AUNMD9gSOPvHCJqXlBxDfQJk")
	fmt.Println(data, err)
}

func TestCustom(t *testing.T) {
	data := &testStruct{"Lorin", 22}
	secret := []byte("secret")

	token, err := jwt.CustomGen(data, time.Minute, secret)
	fmt.Println("token:", token, "err:", err)

	verify := jwt.Verify(token, secret)
	fmt.Println("verify:", verify)

	fmt.Println(jwt.Data(token))
}
