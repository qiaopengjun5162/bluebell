package redis

import (
	"bluebell/setting"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 声明一个全局的 rdb 变量
var (
	client *redis.Client
)

// 初始化连接
func Init(cfg *setting.RedisConfig) (err error) {
	// NewClient将客户端返回给Options指定的Redis Server。
	// Options保留设置以建立redis连接。
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // 没有密码，默认值
		DB:       cfg.DB,       // 默认DB 0 连接到服务器后要选择的数据库。
		PoolSize: cfg.PoolSize, // 最大套接字连接数。 默认情况下，每个可用CPU有10个连接，由runtime.GOMAXPROCS报告。
	})

	// Background返回一个非空的Context。它永远不会被取消，没有值，也没有截止日期。
	// 它通常由main函数、初始化和测试使用，并作为传入请求的顶级上下文
	ctx := context.Background()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
