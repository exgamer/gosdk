package repositories

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
	}
}

// RedisRepository - репозиторий для работы с кешем
type RedisRepository struct {
	redisClient *redis.Client
}

// FlushAll - очистить кеш
func (repo *RedisRepository) FlushAll() error {
	err := repo.redisClient.FlushAll(context.Background()).Err()

	if err != nil {

		return err
	}

	return nil
}
