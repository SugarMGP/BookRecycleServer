package redis

import (
	"bookrecycle-server/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/zjutjh/WeJH-SDK/redisHelper"
	"go.uber.org/zap"
)

// GlobalClient 全局 Redis 客户端实例
var GlobalClient *redis.Client

// Init 初始化 Redis 客户端和配置信息
func Init() {
	redisInfo := redisHelper.InfoConfig{
		Host:     config.Config.GetString("redis.host"),
		Port:     config.Config.GetString("redis.port"),
		DB:       config.Config.GetInt("redis.db"),
		Password: config.Config.GetString("redis.pass"),
	}
	GlobalClient = redisHelper.Init(&redisInfo)
	zap.L().Info("Redis initialized")
}
