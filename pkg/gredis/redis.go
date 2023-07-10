package gredis

import (
	"encoding/json"
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
	return nil
}

// 设置 KEY
func Set(key string, data interface{}, time int) error {
	// 在连接池中获取一个活跃连接
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 向 Redis 服务器发送命令并返回收到的答复
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil
}

// 判断 KEY 是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("GET", key))
	if err != nil {
		return false
	}

	return exists
}

// 获取某个 KEY
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	// 将命令返回转为 Bytes
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// 删除某个 KEY
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	// 将命令返回转为布尔值
	return redis.Bool(conn.Do("DEL", key))
}

// 正则删除某一批 KEY
func likeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	// 将命令返回转为 []string
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	// 循环删除
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
