package redis

import (
	"context"
	"fmt"
	"globe/internal/repository/entities/refreshToken"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	SetRefreshToken(
		ctx context.Context,
		refreshToken *refreshToken.RefreshToken,
		expiratioin time.Duration,
	) error
	GetRefreshTokenByID(
		ctx context.Context,
		id *uint32,
		refreshToken *refreshToken.RefreshToken,
	) error 
}

type repository struct {
	client *redis.Client
}

func InitRepository(client *redis.Client) Repository {
	return &repository{client: client}
}

func (r *repository) SetRefreshToken(
	ctx context.Context,
	refreshToken *refreshToken.RefreshToken,
	expiratioin time.Duration,
) error {
	return r.client.Set(
		ctx,
		strconv.FormatUint(uint64(refreshToken.ID), 10),
		fmt.Sprintf("%s_%v", refreshToken.Token, refreshToken.Expired),
		expiratioin,
	).Err()
}

func (c *repository) GetRefreshTokenByID(
	ctx context.Context,
	id *uint32,
	refreshToken *refreshToken.RefreshToken,
) error {
	refreshTokenExpiration, err := c.client.Get(ctx, strconv.FormatUint(uint64(*id), 10)).Result()
	
	splitRefreshTokenExpiration := strings.Split(refreshTokenExpiration, "_")
	expired, err := time.Parse(splitRefreshTokenExpiration[1], splitRefreshTokenExpiration[1])

	refreshToken.Token = splitRefreshTokenExpiration[0]
	refreshToken.Expired = expired
	return err
}