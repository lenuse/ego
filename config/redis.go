package config

import (
	"github.com/ebar-go/ego/helper"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// RedisConfig redis配置
type RedisConfig struct {
	// 地址
	Host string

	// 端口号
	Port int

	// 密码
	Auth string

	// 连接池大小,默认100个连接
	PoolSize int

	// 最大尝试次数,默认3次
	MaxRetries int

	// 是否为默认连接
	Default bool

	// 超时, 默认10s
	IdleTimeout time.Duration
}

const (
	redisDefaultPort        = 6379
	redisDefaultPoolSize    = 100
	redisDefaultMaxRetries  = 3
	redisDefaultIdleTimeout = 10 * time.Second
)

// complete 自动补全
func (conf *RedisConfig) complete() {
	conf.Port = helper.DefaultInt(conf.Port, redisDefaultPort)
	conf.PoolSize = helper.DefaultInt(conf.PoolSize, redisDefaultPoolSize)
	conf.MaxRetries = helper.DefaultInt(conf.MaxRetries, redisDefaultMaxRetries)

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = redisDefaultIdleTimeout
	}
}

// Options 转换为options
func (conf *RedisConfig) Options() *redis.Options {
	conf.complete()

	address := net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))

	return &redis.Options{
		Addr:        address,
		Password:    conf.Auth,
		PoolSize:    conf.PoolSize,    // Redis连接池大小
		MaxRetries:  conf.MaxRetries,  // 最大重试次数
		IdleTimeout: conf.IdleTimeout, // 空闲链接超时时间
	}
}
