package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/zenfire-cn/webkit/cache"
	"time"
)

func main() {
	cache.Init(cache.DefaultPool("127.0.0.1:6379", 0))

	strings, err := redis.Strings(cache.Get("keys", "*"))
	fmt.Println(strings, err)

	// Bool
	cache.Commit("set", "testBool", true)
	fmt.Println(redis.Bool(cache.Get("get", "testBool")))

	// Int
	cache.Commit("set", "testInt", 503)
	fmt.Println(redis.Int(cache.Get("get", "testInt")))

	// String
	cache.Commit("set", "testString", "string")
	fmt.Println(redis.String(cache.Get("get", "testString")))

	// Hash
	cache.Commit("hset", "testHash", "hKey", "hValue")
	fmt.Println(redis.StringMap(cache.Get("hgetall", "testHash")))
	fmt.Println(redis.String(cache.Get("hget", "testHash", "hKey")))

	// 删除key
	cache.Commit("DEL", "testBool")

	// 判断key是否存在
	fmt.Println(redis.Bool(cache.Get("EXISTS", "testBool")))

	// 设置key过期时间
	cache.Commit("SET", "testExist", true, "EX", "1")
	time.Sleep(time.Second)
	fmt.Println(redis.Bool(cache.Get("EXISTS", "testExist")))

	// 对已有的key设置过期时间
	cache.Commit("EXPIRE", "testInt", 1)
	time.Sleep(time.Second)
	fmt.Println(redis.Bool(cache.Get("EXISTS", "testInt")))
}
