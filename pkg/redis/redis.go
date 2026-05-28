package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/log"
)

type RedisInterface interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type redisClient struct {
	rdb *redis.Client
}

func NewRedisClient() RedisInterface {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.AppEnv.RedisHost + ":" + env.AppEnv.RedisPort,
		Password: env.AppEnv.RedisPassword,
		DB:       0,
	})

	return &redisClient{rdb}
}

func (r *redisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	exp time.Duration,
) error {
	err := r.rdb.Set(ctx, key, value, exp).Err()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[REDIS][Set] faield to set key")

		return err
	}

	return nil
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[REDIS][Get] faield to get key")

		return "", err
	}

	return val, nil
}

func (r *redisClient) Delete(ctx context.Context, key string) error {
	err := r.rdb.Del(ctx, key).Err()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[REDIS][Get] faield to get key")

		return err
	}

	return nil
}
