package redissrorage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dorik33/TgBot/internal/config"
	"github.com/dorik33/TgBot/internal/repository"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		Password: cfg.RedisConfig.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}
	return rdb, nil
}

type Redis struct {
	db  *redis.Client
	ttl time.Duration
}

func NewRedis(db *redis.Client, ttl time.Duration) repository.RedisRepository {
	return &Redis{
		db:  db,
		ttl: ttl,
	}
}

func (r *Redis) SetCryptoPrice(key string, value string) error {
	err := r.db.Set(context.Background(), key, value, r.ttl).Err()
	if err != nil {
		log.Printf("failed to set data, error: %s\n", err.Error())
		return err
	}
	return nil
}

func (r *Redis) GetCryptoPrice(key string) (string, error) {
	price, err := r.db.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("failed to get value, error: %v\n", err)
	}
	return price, err
}

func (r *Redis) Close() {
	r.db.Close()
}
