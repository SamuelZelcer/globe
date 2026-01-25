package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) UpdatePassword(
	ctx context.Context,
	request *dtos.UpdatePasswordRequest,
	token string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Passowrd == "" ||
	request.NewPassowrd == "" ||
	len(request.NewPassowrd) < 8 ||
	len(request.NewPassowrd) > 120 ||
	len(request.Passowrd) < 8 ||
	len(request.Passowrd) > 120 {
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

	// find user
	var user entities.User
	if err := s.userRepository.FindByID(claims.UserID, &user); err != nil {
		return nil, errors.New("Couldn't find user")
	}

	// authenticate user
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(request.Passowrd)); err != nil {
		return nil, errors.New("Wrong password, authentication failed")
	}

	// hash new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassowrd), 10)
	if err != nil {
		return nil, errors.New("Couldn't generate hashed passowd")
	}

	// send emsil with verification code
	code := strconv.Itoa(rand.Intn(900000) + 100000)
	go s.email.SendVerificationCode(code, user.Email)

	// parse new password and verification code
	newPasswordAndCodeJSON, err := json.Marshal(
		&dtos.NewPasswordAndVerificationCode{
			Code: code,
			NewPassowrd: string(hashedNewPassword),
		},
	)

	// save new password and verification code to cache
	if err := s.redis.SET(
		ctx,
		fmt.Sprintf("newpassword:%d", claims.UserID),
		string(newPasswordAndCodeJSON),
		time.Minute*15,
	); err != nil {
		return nil, errors.New("Coludn't save new password and verifiation code to cache")
	}
	return &tokens, nil
}