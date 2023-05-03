package services

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