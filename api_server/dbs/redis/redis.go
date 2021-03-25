package redis

import (
	"swc/logger"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	maxIdle     = 10               // 空闲时维持的最大连接数
	maxActive   = 200              // 维持的最大连接数
	idleTimeout = 60 * time.Second // 每个连接维持的最长时间
)

type redisPoll struct {
	*redigo.Pool
}

var (
	poll redisPoll // 连接池
)

// 初始化 redis 信息, 可以从别的地方进行配置
func InitRedis(addr, passward string) {
	pool, err := poolInitRedis(addr, passward)
	if err != nil {
		logger.Error.Fatal("redis 连接失败")
	}
	poll.Pool = pool
}

// redis pool
func poolInitRedis(server string, password string) (pool *redigo.Pool, err error) {
	pool = &redigo.Pool{
		MaxIdle:      maxIdle,
		IdleTimeout:  idleTimeout,
		MaxActive:    maxActive,
		Dial:         dial(server, password),
		TestOnBorrow: ping,
	}
	c := pool.Get()
	_, err = c.Do("PING")
	c.Close()
	return
}

// 向连接池提供 reids 地址和 密码
func dial(server string, password string) func() (redigo.Conn, error) {
	return func() (redigo.Conn, error) {
		c, err := redigo.Dial("tcp", server) // 建立tcp连接
		if err != nil {
			return nil, err
		}
		if password != "" { // 是否需要输入密码
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	}
}

// 连通性测试
func ping(c redigo.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}

// Get 返回 redis 的连接对象, 使用完应该释放
func Get() redigo.Conn {
	return poll.Get()
}
