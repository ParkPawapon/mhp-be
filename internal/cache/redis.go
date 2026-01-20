package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/ParkPawapon/mhp-be/internal/config"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
