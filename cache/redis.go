package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

func Init(p *redis.Pool) {
	pool = p
}

func DefaultPool(host string, db int) *redis.Pool {
	return &redis.Pool{
		MaxActive:   10,                               // 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
		MaxIdle:     20,                               // 最大空闲连接数
		IdleTimeout: time.Duration(120) * time.Second, // 空闲连接超时时间
		Wait:        true,                             // 当超过最大连接数 是报错还是等待， true 等待 false 报错
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialPassword(""), redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

func Conn() redis.Conn {
	return pool.Get()
}

func Get(commandName string, args ...interface{}) (interface{}, error) {
	conn := Conn()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func Commit(commandName string, args ...interface{}) {
	conn := Conn()
	defer conn.Close()
	conn.Do(commandName, args...)
	conn.Flush()
}