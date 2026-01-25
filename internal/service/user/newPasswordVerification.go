package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
)

func (s *service) NewPasswordVerification(
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

	// find new password and verification code from cache
	newPasswordAndCodeJSON, err := s.redis.GET(ctx, fmt.Sprintf("newpassword:%d", claims.UserID))
	if err != nil {
		return nil, errors.New("Couldn't find new password in cache")
	}
	
	// parse new password and verification code from JSON
	var newPasswordAndCode dtos.NewPasswordAndVerificationCode 
	if err := json.Unmarshal([]byte(newPasswordAndCodeJSON), &newPasswordAndCode); err != nil {
		return nil, errors.New("Couldn't parse new password from cache")
	}

	// update user in database
	if err := s.userRepository.UpdatePasswordByID(claims.UserID, newPasswordAndCode.NewPassowrd); err != nil {
		return nil, errors.New("Couldn't save updated user to database")
	}

	// delete outdated data from redis
	s.redis.DEL(ctx, fmt.Sprintf("newpassword:%d", claims.UserID))


	return &tokens, nil
}