package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var _ Client = (*client)(nil)

type Client interface {
	Ping() error
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Close() error
}

type client struct {
	client *redis.Client
}

func NewClient(addr string) *client {
	return &client{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *client) Ping() error {
	return c.client.Ping().Err()
}

func (c *client) Get(key string) (string, error) {
	res, err := c.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (c *client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(key, value, expiration).Err()
}

func (c *client) Close() error {
	return c.client.Close()
}
