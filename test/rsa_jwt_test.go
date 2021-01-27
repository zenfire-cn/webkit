package test

import (
	"fmt"
	"github.com/zenfire-cn/webkit/jwt"
	"testing"
	"time"
)

/**
 * @description: 可以先通过以下命令生成公钥和私钥
 *  openssl genrsa -out private.key 2048
 *  openssl rsa -in private.key -pubout > public.key
 */
func TestRsaGen(t *testing.T) {
	t.Run("TestRsaGen", func(t *testing.T) {
		jwt.RsaInit(time.Minute * 2, "public.key", "private.key")

		token, err := jwt.RsaGen(map[string]string{"id": "1", "phone": "10086"})
		fmt.Println(token, err)
	})
}

func TestRsaParse(t *testing.T) {
	t.Run("TestRsaParse", func(t *testing.T) {
		jwt.RsaInit(time.Minute * 2, "public.key", "private.key")

		res, err := jwt.RsaParse("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7ImlkIjoiMSIsInBob25lIjoiMTAwODYifSwiZXhwIjoxNjExNzE5MjgxfQ.m1Y3krLWMJ2Cba_c74s5Cr8Ko8DTeOiQliplwurpF9ElZnQDe1_fuy95umrHUIdEQkGG-lVcq3WYchj4H97papLDkayoFjipO-D4fX8VNkmFoJMZphgsLL-5hWdSZxZy40XjyvKx-2i0YluDAwB851wmQJnOR10J_1qcR5Q0xZA")
		fmt.Println(res, err)
	})
}
