package redis

import (
	"OnlineJudge/setting"
	"context"
	"fmt"
	"github.com/go-redis/redis"
)

var (
	Client *redis.Client
	Nil    = redis.Nil
	Ctx    = context.Background()
)

func Init(cfg *setting.RedisConfig) (err error) {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = Client.Ping().Result()
	return
}

func Close() {
	_ = Client.Close()
}
