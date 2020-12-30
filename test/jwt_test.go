package test

import (
	"fmt"
	"github.com/zenfire-cn/webkit/jwt"
	"testing"
	"time"
)

func init() {
	jwt.Init(time.Hour*2, "mySecret")
}

type testStruct struct {
	Name string
	Age int
}

func TestGen(t *testing.T) {
	t.Run("TestGen", func(t *testing.T) {
		token, err := jwt.Gen(&testStruct{"Lorin", 22})
		fmt.Println(token, err)
	})
}

func TestParse(t *testing.T) {
	t.Run("TestParse", func(t *testing.T) {
		res, err := jwt.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7Ik5hbWUiOiJMb3JpbiIsIkFnZSI6MjJ9LCJleHAiOjE2MDkyMzM5MTN9.56n-pHdWPpR8I9INQLCyyvVzVUf43y9jUbO01X6X1qw")
		fmt.Println(res, err)
	})
}
