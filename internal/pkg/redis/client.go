package redis

import (
	"context"
	"time"

	"ride-sharing/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	cli *redis.Client
}

func New(cfg *config.Config) *Client {
	return &Client{
		cli: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

func (c *Client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return c.cli.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.cli.Close()
}
