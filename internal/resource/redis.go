package resource

import (
	"fmt"
	"product-catalog-service/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis(appConfig config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", appConfig.Redis.Host, appConfig.Redis.Port),
		Password: appConfig.Redis.Password,
	})
	return redisClient
}
