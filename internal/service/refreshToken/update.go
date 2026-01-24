package refreshTokenService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	JWT "globe/internal/service/jwt"
	"strings"
	"time"
)

func (s *service) Update(
	ctx context.Context,
	providedRefreshToken string,
	providedAccessToken string,
	tokens *dtos.AuthenticationTokens,
) (*JWT.UserClaims, error) {
	// validate pair of tokens
	if providedRefreshToken == "" ||
	providedAccessToken == "" {
		return nil, errors.New("Bad request")
	}

	// validate access token signature and extract claims
	claims, err := s.jwtManager.ValidateWithoutExpiration(providedAccessToken)
	if err != nil {
		return nil, errors.New("Invalid access token")
	}

	// find refresh token
	var refreshToken entities.RefreshToken
	value, err := s.redis.GET(ctx, fmt.Sprintf("refreshtoken:%d", claims.UserID))
	if err != nil {
		if err := s.refreshTokenRepository.FindByID(claims.UserID, &refreshToken); err != nil {
			return nil, errors.New("Couldn't find refresh token")
		}
	}
	// parse value
	splitValue := strings.Split(value, "_")
	expiration, err := time.Parse(time.RFC3339, splitValue[1])
	if err != nil {
		return nil, errors.New("Couldn't parse expiration bact to time.Time")
	}
	refreshToken.Token = splitValue[0]
	refreshToken.Expired = expiration

	// validate provided refreshToken
	if refreshToken.Token != providedRefreshToken || refreshToken.Expired.Before(time.Now()) {
		return nil, errors.New("Token is invalid or expired")
	}

	// create new refresh token
	newRefreshToken, err := s.Create(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("Couldn't create new refresh token")
	}

	// new access token
	newAccessToken, err := s.jwtManager.Create(claims.UserID, claims.Subject, time.Minute*5)
	if err != nil {
		return nil, errors.New("Couldn't generate new access token")
	}
	tokens.RefreshToken = newRefreshToken
	tokens.AccessToken = newAccessToken
	return claims, nil
}