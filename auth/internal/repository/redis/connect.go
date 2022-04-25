package redis

import (
	"context"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	c *redis.Client
}

func Connect() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        utils.GetConfig().Redis.Address,
		Password:    utils.GetConfig().Redis.Password,
		ReadTimeout: utils.GetConfig().Redis.RTimeout * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), utils.GetConfig().Redis.RTimeout*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewRedis(conn *redis.Client) *Redis {
	return &Redis{conn}
}
