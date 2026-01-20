package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	SET(
		ctx context.Context,
		key string,
		value string,
		expiration time.Duration,
	) error
	GET(
		ctx context.Context,
		key string,
	) (string, error)
	DEL(
		ctx context.Context,
		key string,
	) error
	EXPIRE(
		ctx context.Context,
		key string,
		expiration time.Duration,
	) error
}

type cache struct {
	client *redis.Client
}

func InitRepository(client *redis.Client) Cache {
	return &cache{client: client}
}

func (c *cache) SET(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	return c.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}

func (c *cache) GET(
	ctx context.Context,
	key string,
) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *cache) DEL(
	ctx context.Context,
	key string,
) error {
	return c.client.Del(
		ctx,
		key,
	).Err()
}

func (c *cache) EXPIRE(
	ctx context.Context,
	key string,
	expiration time.Duration,
) error {
	return c.client.Expire(
		ctx,
		key,
		expiration,
	).Err()
}