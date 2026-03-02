package check

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCheck struct {
	URL    string
	pinger func(ctx context.Context, url string) error
}

func (c *RedisCheck) Name() string {
	return "Redis reachable"
}

func (c *RedisCheck) Run(ctx context.Context) Result {
	ping := c.pinger
	if ping == nil {
		ping = func(ctx context.Context, url string) error {
			opt, err := redis.ParseURL(url)
			if err != nil {
				return err
			}
			client := redis.NewClient(opt)
			defer client.Close()
			return client.Ping(ctx).Err()
		}
	}

	if err := ping(ctx, c.URL); err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: "cannot reach Redis",
			Fix:     "make sure Redis is running and REDIS_URL is correct",
		}
	}
	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "Redis is reachable",
	}
}
