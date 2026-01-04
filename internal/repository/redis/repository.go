package redis

import (
	"context"
	"globe/internal/repository/entities/refreshToken"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	SaveRefreshToken(
		ctx context.Context,
		refreshToken *refreshToken.RefreshToken,
		expiratioin time.Duration,
	) error
}

type repository struct {
	client *redis.Client
}

func InitRepository(client *redis.Client) Repository {
	return &repository{client: client}
}

func (r *repository) SaveRefreshToken(
	ctx context.Context,
	refreshToken *refreshToken.RefreshToken,
	expiratioin time.Duration,
) error {
	return r.client.Set(
		ctx,
		strconv.FormatUint(uint64(refreshToken.ID), 10),
		refreshToken.Token,
		expiratioin,
	).Err()
}