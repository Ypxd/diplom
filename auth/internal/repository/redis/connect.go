package redis

import (
	"context"
	"encoding/json"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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

func (r *Redis) UpdateUserTags(t string, userID string) (map[int64]int64, error) {
	var val [10]string
	keys, err := r.c.Keys(context.Background(), userID).Result()
	if err != nil {
		return nil, errors.WithMessage(err, "get keys")
	}
	if len(keys) == 0 {
		val[0] = t
		vByte, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		err = r.c.Set(context.Background(), userID, vByte, time.Duration(10*time.Hour)).Err()
		if err != nil {
			return nil, errors.WithMessage(err, "set keys")
		}
		return nil, nil
	}

	s := r.c.Get(context.Background(), keys[0]).Val()
	err = json.Unmarshal([]byte(s), &val)
	if err != nil {
		return nil, err
	}
	for i, v := range val {
		if v == "" {
			val[i] = t
			vByte, err := json.Marshal(val)
			if err != nil {
				return nil, err
			}
			err = r.c.Set(context.Background(), keys[0], vByte, time.Duration(10*time.Hour)).Err()
			if err != nil {
				return nil, errors.WithMessage(err, "set keys")
			}
			return nil, nil
		}
	}

	for i := 0; i < len(val)-1; i++ {
		val[i] = val[i+1]
	}
	val[9] = t
	vByte, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	err = r.c.Set(context.Background(), keys[0], vByte, time.Duration(10*time.Hour)).Err()
	if err != nil {
		return nil, errors.WithMessage(err, "set keys")
	}

	res := make(map[int64]int64, 0)
	for _, v := range val {
		if v != "" {
			tag := strings.Split(v, ";")
			for j := range tag {
				singleTag, err := strconv.ParseInt(tag[j], 10, 64)
				if err != nil {
					return nil, err
				}
				res[singleTag]++
			}
		}
	}

	return res, nil
}

func NewRedis(conn *redis.Client) *Redis {
	return &Redis{conn}
}
