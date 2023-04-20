package services

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

type RedisService struct {
	redisPool *redis.Pool
}

// 初始化RedisService
func NewRedisService() *RedisService {
	redis_host := beego.AppConfig.String("redis_host")
	redis_port := beego.AppConfig.String("redis_port")
	redis_password := beego.AppConfig.String("redis_password")
	redis_db := beego.AppConfig.String("redis_db")
	redis_addr := fmt.Sprintf("%s:%s", redis_host, redis_port)

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redis_addr)
			if err != nil {
				return nil, err
			}
			if redis_password != "" {
				if _, err := c.Do("AUTH", redis_password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if redis_db != "" {
				if _, err := c.Do("SELECT", redis_db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}

	return &RedisService{redisPool: redisPool}
}

// 关闭Redis连接池
func (s *RedisService) Close() error {
	return s.redisPool.Close()
}

// 从Redis连接池获取连接
func (s *RedisService) GetConn() redis.Conn {
	return s.redisPool.Get()
}

// 执行Redis命令
func (s *RedisService) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := s.GetConn()
	defer conn.Close()

	return conn.Do(commandName, args...)
}
