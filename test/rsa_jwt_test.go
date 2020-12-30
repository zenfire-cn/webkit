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
func init() {
	jwt.RsaInit(time.Hour * 2, "public.key", "private.key")
}

func TestRsaGen(t *testing.T) {
	t.Run("TestRsaGen", func(t *testing.T) {
		token, err := jwt.RsaGen(map[string]string{"id": "1", "phone": "10086"})
		fmt.Println(token, err)
	})
}

func TestRsaParse(t *testing.T) {
	t.Run("TestRsaParse", func(t *testing.T) {
		res, err := jwt.RsaParse("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7ImlkIjoiMSIsInBob25lIjoiMTAwODYifSwiZXhwIjoxNjA5MjMzOTUyfQ.mhUt85rs_7r5hHmXy4rHTyTC-7yYoiU6YWIr91DmGKT2cwj9qnPdBMviT9SXR2j_FzqoM-WvhOPh-x8oS3_5wmOj6j8QbsstM4qrg2BlkSRt6c5P6UCY2ybCVfLz0ZQyChkJoUz0abV4PVijuvCyK7IrDrPQCm_9U_BzFgptSU8")
		fmt.Println(res, err)
	})
}
