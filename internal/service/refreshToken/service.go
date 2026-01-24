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
	Update(
		ctx context.Context,
		providedRefreshToken string,
		providedAccessToken string,
		tokens *dtos.AuthenticationTokens,
	) (*JWT.UserClaims, error)
	Create(ctx context.Context, userID uint64) (string, error)
}

type service struct {
	refreshTokenRepository refreshToken.Repository
	userRepository user.Repository
	jwtManager JWT.Manager
	redis redis.Cache
}

func Init(
	refreshTokenRepository refreshToken.Repository,
	userRepository user.Repository,
	jwtManager JWT.Manager,
	redis redis.Cache,
) Service {
	return &service{
		refreshTokenRepository: refreshTokenRepository,
		userRepository: userRepository,
		jwtManager: jwtManager,
		redis: redis,
	}
}