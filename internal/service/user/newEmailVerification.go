package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"time"
)

func (s *service) NewEmailVerification(
	ctx context.Context,
	request *dtos.VerificationRequest,
	token string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Code == "" ||
	len(request.Code) != 6 {
		return nil, errors.New("Invalid request")
	}

	// validate token
	var tokens dtos.AuthenticationTokens
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == "" {
			return nil, errors.New("Invalid refresh token")
		}

		// update authentication tokens
		claims, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, &tokens)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	// find new email and verification code from cache
	newEmailAndCodeJSON, err := s.redis.GET(ctx, fmt.Sprintf("newemail:%d", claims.UserID))
	if err != nil {
		return nil, errors.New("Couldn't find new email in cache")
	}
	
	// parse new email and verification code from JSON
	var newEmailAndCode dtos.NewEmailAndVerificationCode 
	if err := json.Unmarshal([]byte(newEmailAndCodeJSON), &newEmailAndCode); err != nil {
		return nil, errors.New("Couldn't parse new email from cache")
	}

	// update user in database
	if err := s.userRepository.UpdateEmailByID(claims.UserID, newEmailAndCode.Email); err != nil {
		return nil, errors.New("Couldn't save updated user to database")
	}

	// update access token
	newAccessToken, err := s.jwtManager.Create(claims.UserID, newEmailAndCode.Email, time.Minute*5)
	if err != nil {
		return nil, errors.New("Couldn't create new access token")
	}
	tokens.AccessToken = newAccessToken

	// delete outdated data from redis
	s.redis.DEL(ctx, fmt.Sprintf("newemail:%d", claims.UserID))

	return &tokens, nil
}