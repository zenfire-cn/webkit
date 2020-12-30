package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zenfire-cn/commkit/utility"
	"github.com/zenfire-cn/webkit/rest"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Jwt struct {
	Data interface{}
	jwt.StandardClaims
}

type option struct {
	expire     time.Duration
	secret     []byte
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	rsa        bool
}

var (
	o        *option
	callback func(jwtData map[string]interface{}, validations ...interface{}) (bool, *rest.Rest)
)

func Init(expire time.Duration, secret string) {
	o = &option{
		expire: expire,
		secret: []byte(secret),
	}
}

func RsaInit(expire time.Duration, publicKey, privateKey string) {
	pub, _ := jwt.ParseRSAPublicKeyFromPEM(readKeyFile(publicKey))
	pri, _ := jwt.ParseRSAPrivateKeyFromPEM(readKeyFile(privateKey))
	o = &option{
		expire:     expire,
		publicKey:  pub,
		privateKey: pri,
		rsa:        true,
	}
}

func AuthCallback(f func(jwtData map[string]interface{}, validations ...interface{}) (bool, *rest.Rest)) {
	callback = f
}

func readKeyFile(path string) []byte {
	filePath := utility.FindConfigFile(path)
	if filePath == "" {
		log.Fatal("could not found public key file")
	}
	keyByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	return keyByte
}

func Gen(data interface{}) (string, error) {
	j := Jwt{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(o.expire).Unix(), // 过期时间
		},
	}
	// 使用jwt库中已有的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(o.secret)
}

func Parse(tokenStr string) (map[string]interface{}, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenStr, &Jwt{}, func(token *jwt.Token) (i interface{}, err error) {
		return o.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Jwt); ok && token.Valid { // 校验token
		res := claims.Data.(map[string]interface{})
		res["expiresAt"] = claims.ExpiresAt
		return res, nil
	}
	return nil, errors.New("invalid token")
}

func RsaGen(data interface{}) (string, error) {
	j := Jwt{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(o.expire).Unix(), // 过期时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, j)
	return token.SignedString(o.privateKey)
}

func RsaParse(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Jwt{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
		}
		return o.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Jwt); ok && token.Valid { // 校验token
		res := claims.Data.(map[string]interface{})
		res["expiresAt"] = claims.ExpiresAt
		return res, nil
	}
	return nil, errors.New("invalid token")
}

func Auth(validations ...interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			rest.Error(c, "请求头中auth为空")
			c.Abort()
			return
		}
		// "Bearer tokenString..." 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			rest.Error(c, "请求头中auth格式有误")
			c.Abort()
			return
		}

		var (
			jwtData map[string]interface{}
			err     error
		)

		if o.rsa {
			jwtData, err = RsaParse(parts[1])
		} else {
			jwtData, err = Parse(parts[1])
		}

		if err != nil {
			rest.Error(c, "无效的token")
			c.Abort()
			return
		}

		if callback != nil {
			if ok, response := callback(jwtData, validations...); !ok {
				c.JSON(http.StatusOK, response)
				c.Abort()
				return
			}
		}

		// 将当前请求的token信息保存到请求的上下文里，后续的处理函数可以用过 c.Get("jwtData") 来获取当前请求的用户信息
		c.Set("jwtData", jwtData)
		c.Next()
	}
}
