package api

import (
	"context"

	goRedis "github.com/redis/go-redis/v9"
)

var redis goRedis.Client = *goRedis.NewClient(
    &goRedis.Options{
        Addr: REDIS_HOST+":"+REDIS_PORT,
        Password: REDIS_PASSWORD,
        DB: REDIS_DB,
    },
)
var redisContext context.Context = context.Background()

// Returns whether a redis key exists or not
func redisKeyExists(key string) (bool, error) {
    exists, err := redis.Exists(redisContext, key).Result()
    if err != nil {
        return false, err
    }
    if exists == 1 {
        return true, nil
    }
    return false, nil
}