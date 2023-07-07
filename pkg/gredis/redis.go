package gredis

import (
	"go-gin-blog-api/pkg/setting"
	"time"

	// Redigo是Redis数据库的Go客户端
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// 启动
func Setup() error {
	// 采用 Redis 连接池
	RedisConn = &redis.Pool{
		// 最大空闲连接数
		MaxIdle: setting.RedisSetting.MaxIdle,
		// 在给定时间内，允许分配的最大连接数
		MaxActive: setting.RedisSetting.MaxActive,
		// 在给定时间内，保持空闲状态的时间，若到达时间限制则关闭连接
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		// 提供创建和配置应用程序连接的一个函数
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			// 采用密码 auth 指令
			if setting.RedisSetting.Password != "" {
				if _, err := conn.Do("AUTH", setting.RedisSetting.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		// 健康检测
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
