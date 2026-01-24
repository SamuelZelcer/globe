package userService

import (
	"context"
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignIn(request *dtos.SignInRequest, ctx context.Context) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Email == "" ||
	request.Password == "" ||
	len(request.Password) < 8 ||
	len(request.Password) > 120 {
	    return nil, errors.New("Bad request")
	}

	// validate email
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, errors.New("Invalid email")
	}

	// find user
	var user entities.User
	if err := s.userRepository.FindByEmail(request.Email, &user);
	err != nil || user.Email == "" {
		return nil, errors.New("Couldn't find user by email")
	}


	// compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("Invalid password provided")
	}

	// create refresh token
	refreshToken, err := s.refreshTokenService.Create(ctx, user.ID)
	if err != nil {
		return nil, errors.New("Couldn't create refresh token")
	}

	// access token
	accessToken, err := s.jwtManager.Create(user.ID, user.Email, time.Minute*5)
	if err != nil {
		return nil, errors.New("Couldn't generate access token")
	}

	return &dtos.AuthenticationTokens{
		RefreshToken: refreshToken,
		AccessToken: accessToken,
	}, nil
}