package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"online_judge/setting"
	"time"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func Init(cfg *setting.RedisConfig) (err error) {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = Client.Ping(Ctx).Result()

	interval := time.Duration(cfg.PersistenceInterval) * time.Hour
	go Durability(interval)
	return
}

func Durability(interval time.Duration) {
	// RDB持久化
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := Client.BgSave(Ctx).Err()
			if err != nil {
				zap.L().Error("redis save fail", zap.Error(err))
			}
		}
	}
}

func Close() {
	_ = Client.Close()
}
