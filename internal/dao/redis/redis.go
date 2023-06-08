package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RDB *redis.Client

func Init() error {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			viper.GetString("redis.host"),
			viper.GetString("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolFIFO: false,
		// PoolTimeout 代表如果连接池所有连接都在使用中，等待获取连接时间，超时将返回错误
		// 默认是 1秒+ReadTimeout
		PoolTimeout: time.Duration(1 + 13),
	})

	var ctx context.Context
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	RDB = rdb
	return nil
}
