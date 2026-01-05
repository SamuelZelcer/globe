package refreshTokenService

import (
	"context"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/refreshToken"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	JWT "globe/internal/service/jwt"
)

type Service interface {
	Update(ctx context.Context, refreshToken *string, accessToken *string) (*dtos.AuthenticationTokens, error)
}

type service struct {
	refreshTokenRepository refreshToken.Repository
	userRepository user.Repository
	jwtManager JWT.Manager
	redis redis.Repository
}

func Init(
	refreshTokenRepository refreshToken.Repository,
	userRepository user.Repository,
	jwtManager JWT.Manager,
	redis redis.Repository,
) Service {
	return &service{
		refreshTokenRepository: refreshTokenRepository,
		userRepository: userRepository,
		jwtManager: jwtManager,
		redis: redis,
	}
}