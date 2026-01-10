package refreshTokenService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"time"

	"github.com/google/uuid"
)

func (s *service) Update(
	ctx context.Context,
	providedRefreshToken *string,
	providedAccessToken *string,
) (*dtos.AuthenticationTokens, error) {
	// validate pair of tokens
	if providedRefreshToken == nil ||
	providedAccessToken == nil {
		return nil, errors.New("Bad request")
	}

	// validate access token signature and extract claims
	claims, err := s.jwtManager.ValidateWithoutExpiration(providedAccessToken)
	if err != nil {
		return nil, errors.New("Invalid access token")
	}

	// find refresh token
	var refreshToken entities.RefreshToken
	if err := s.redis.GetRefreshTokenByID(ctx, &claims.UserID, &refreshToken); err != nil {
		if err := s.refreshTokenRepository.FindByID(&claims.UserID, &refreshToken); err != nil {
			return nil, errors.New("Couldn't find refresh token")
		}
	}

	// validate provided refreshToken
	if refreshToken.Token != *providedRefreshToken || refreshToken.Expired.Before(time.Now()) {
		fmt.Printf("%s   ---   %s\n", refreshToken.Token, *providedRefreshToken)
		return nil, errors.New("Token is invalid or expired")
	}

	// new refresh token
	newRefreshToken := &entities.RefreshToken{
		ID: claims.UserID,
		Token: uuid.NewString(),
		Expired: time.Now().Add(time.Hour*168),
	}

	// save new refresh token
	if err := s.redis.SetRefreshToken(ctx, newRefreshToken, time.Hour*24); err != nil {
		return nil, errors.New("Couldn't save refresh token to redis")
	}
	if err := s.refreshTokenRepository.Save(newRefreshToken); err != nil {
		return nil, errors.New("Couldn't save refresh token to DB")
	}

	// new access token
	newAccessToken, err := s.jwtManager.Create(&claims.UserID, &claims.Subject, time.Minute*5)
	if err != nil {
		return nil, errors.New("Couldn't generate new access token")
	}
	return &dtos.AuthenticationTokens{
		RefreshToken: &newRefreshToken.Token,
		AccessToken: newAccessToken,
	}, nil
}