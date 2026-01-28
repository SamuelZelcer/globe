package userUpdateService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"math/rand"
	"net/mail"
	"strconv"
	"time"
)

func (s *userUpdateService) UpdateEmail(
	ctx context.Context,
	request *dtos.UpdateEmailRequest,
	token string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Email == "" {
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

	// validate email
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, errors.New("Bad request")
	}

	// is email already in use
	isEmailAlreadyInUse, err := s.userRepository.IsEmailAlreadyInUse(request.Email)
	if err != nil {
		return nil, errors.New("Couldn't check is email already in use")
	} else if (isEmailAlreadyInUse) {
		return nil, errors.New("Email already in use")
	}

	// verification code
	code := strconv.Itoa(rand.Intn(900000) + 100000)

	// send verification code to new email
	go s.email.SendVerificationCode(code, request.Email)

	// parse new email and verification code
	newEmailAndCodeJSON, err := json.Marshal(
		dtos.NewEmailAndVerificationCode{
			Email: request.Email,
			Code: code,
		},
	)

	// save new email and verification code to cache
	if err := s.redis.SET(
		ctx,
		fmt.Sprintf("newemail:%d", claims.UserID),
		string(newEmailAndCodeJSON),
		time.Minute*15,
	); err != nil {
		return nil, errors.New("Couldn't save new email and verification code to cache")
	}
	
	return &dtos.AuthenticationTokens{
		RefreshToken: tokens.RefreshToken,
		AccessToken: tokens.AccessToken,
	}, nil
}