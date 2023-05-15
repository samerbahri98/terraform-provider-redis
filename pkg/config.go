package pkg

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Client struct {
	client *redis.Client
}

func (c *Config) Client() (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})
	error := rdb.Ping(ctx).Err()
	return &Client{client: rdb}, error
}
