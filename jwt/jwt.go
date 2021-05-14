package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
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
	Iss string `json:"iss"`
}

type option struct {
	expire      time.Duration
	secret      []byte
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
	ExpiredCode int
	ExpiredMsg  string
	rsa         bool
	Iss         string
}

var (
	o        *option
	callback func(jwtData map[string]interface{}, validations ...interface{}) (bool, *rest.Rest)
)

func Init(expire time.Duration, secret string) {
	o = &option{
		expire:      expire,
		secret:      []byte(secret),
		ExpiredCode: 510,
		ExpiredMsg:  "invalid token",
	}
}

func SetIss(iss string) {
	o.Iss = iss
}

func SetExpiredCodeAndMsg(code int, msg string) {
	o.ExpiredCode = code
	o.ExpiredMsg = msg
}

func RsaInit(expire time.Duration, publicKey, privateKey string) {
	pub, _ := jwt.ParseRSAPublicKeyFromPEM(readKeyFile(publicKey))
	pri, _ := jwt.ParseRSAPrivateKeyFromPEM(readKeyFile(privateKey))
	o = &option{
		expire:     expire,
		ExpiredCode: 510,
		ExpiredMsg:  "invalid token",
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

/**
 * 根据初始化的配置，生成token
 */
func Gen(data interface{}) (string, error) {
	return CustomGen(data, o.expire, o.secret, o.Iss)
}

/**
 * 根据初始化的配置，解析token
 */
func Parse(tokenStr string) (map[string]interface{}, error) {
	return CustomParse(tokenStr, o.secret)
}

/**
 * 生成token
 */
func CustomGen(data interface{}, expire time.Duration, secret []byte, iss string) (string, error) {
	j := Jwt{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(), // 过期时间
		},
		Iss: iss,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	return token.SignedString(secret)
}

/**
 * 校验token有效性，并返回token中的数据
 */
func CustomParse(tokenStr string, secret []byte) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Jwt{}, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Jwt); ok && token.Valid {
		res := claims.Data.(map[string]interface{})
		res["expiresAt"] = claims.ExpiresAt
		return res, nil
	}
	return nil, errors.New("invalid token")
}

/**
 * 只校验token有效性，不返回token中的数据
 */
func Verify(tokenStr string, secret []byte) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err == nil {
		return token.Valid
	}
	return false
}

/**
 * 不校验token有效性，直接获取token payload中的数据
 */
func Data(token string) (map[string]interface{}, error) {
	split := strings.Split(token, ".")
	if len(split) > 1 {
		decodeStr, err := base64.RawStdEncoding.DecodeString(split[1])
		j := &Jwt{}
		json.Unmarshal(decodeStr, j)
		data := j.Data.(map[string]interface{})
		data["expiresAt"] = j.ExpiresAt
		return data, err
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

/**
 * gin路由中间件
 */
func Auth(validations ...interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			rest.Error(c, "请求头缺少Authorization字段")
			c.Abort()
			return
		}
		// "Bearer tokenString..." 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			rest.Error(c, "请求头Authorization字段格式有误")
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
			rest.New(c, o.ExpiredCode, nil, o.ExpiredMsg)
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
