package refreshTokenService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	JWT "globe/internal/service/jwt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *service) Update(
	ctx context.Context,
	providedRefreshToken *string,
	providedAccessToken *string,
	tokens *dtos.AuthenticationTokens,
) (*JWT.UserClaims, error) {
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
	value, err := s.redis.GET(ctx, strconv.FormatUint(uint64(claims.UserID), 10))
	if err != nil {
		if err := s.refreshTokenRepository.FindByID(&claims.UserID, &refreshToken); err != nil {
			return nil, errors.New("Couldn't find refresh token")
		}
	}
	// parse value
	spltValue := strings.Split(value, "_")
	expiration, err := time.Parse(time.RFC3339, spltValue[1])
	if err != nil {
		return nil, errors.New("Couldn't parse expiration bact to time.Time")
	}
	refreshToken.Token = spltValue[0]
	refreshToken.Expired = expiration

	// validate provided refreshToken
	if refreshToken.Token != *providedRefreshToken || refreshToken.Expired.Before(time.Now()) {
		return nil, errors.New("Token is invalid or expired")
	}

	// new refresh token
	newRefreshToken := &entities.RefreshToken{
		ID: claims.UserID,
		Token: uuid.NewString(),
		Expired: time.Now().Add(time.Hour*168),
	}

	// save new refresh token to redis
	if err := s.redis.SET(
		ctx,
		strconv.FormatUint(uint64(claims.UserID), 10),
		fmt.Sprintf("%s_%s", newRefreshToken.Token, newRefreshToken.Expired.Format(time.RFC3339)),
		time.Hour*24,
	); err != nil {
		return nil, errors.New("Couldn't store refresh token in redis")
	}

	// save new refresh token to database
	if err := s.refreshTokenRepository.Save(newRefreshToken); err != nil {
		return nil, errors.New("Couldn't save refresh token to DB")
	}

	// new access token
	newAccessToken, err := s.jwtManager.Create(&claims.UserID, &claims.Subject, time.Minute*5)
	if err != nil {
		return nil, errors.New("Couldn't generate new access token")
	}
	tokens.RefreshToken = &newRefreshToken.Token
	tokens.AccessToken = newAccessToken
	return claims, nil
}