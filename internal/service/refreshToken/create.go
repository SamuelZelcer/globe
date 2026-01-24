package refreshTokenService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/entities"
	"time"

	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, userID uint64) (string, error) {
	// refresh token
	refreshToken := entities.RefreshToken{
		ID: userID,
		Token: uuid.NewString(),
		Expired: time.Now().Add(time.Hour*168),
	}

	// save refresh token to DB
	if err := s.refreshTokenRepository.Save(&refreshToken); err != nil {
		return "", errors.New("Couldn't create refresh token")
	}

	// save refresh token to redis
	if err := s.redis.SET(
		ctx,
		fmt.Sprintf("refreshtoken:%d", refreshToken.ID),
		fmt.Sprintf("%s_%s", refreshToken.Token, refreshToken.Expired.Format(time.RFC3339)),
		time.Hour*24,
	); err != nil {
		return "", errors.New("Couldn't store refresh token in redis")
	}
	return refreshToken.Token, nil
}