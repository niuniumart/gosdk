package redislock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

// Client Redis 客户端.
type Client struct {
	pool *redis.Client
}

func NewClient(address, password string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})

	// 保证能ping通
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redislock ping err is %v", err)
	}

	return &Client{pool: client}
}

// func (c *Client) GetConn(ctx context.Context) *redis.Conn {
// 	return c.pool.Conn()
// }
