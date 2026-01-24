package userService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"regexp"
)

func (s *service) UpdateUsername(
	ctx context.Context,
	request *dtos.UpdateUsernameRequest,
	token string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Username == "" ||
	len(request.Username) < 4 ||
	len(request.Username) > 60 {
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

	// validate username
	matched, err := regexp.MatchString(`[!@#$%^&*()]`, request.Username)
	if err != nil || matched {
		return nil, errors.New("Invalid username")
	}

	// is username already in use
	isUsernameAlreadyInUse, err := s.userRepository.IsUsernameAlreadyInUse(request.Username)
	if err != nil {
		return nil, errors.New("Couldn't check is username already in use")
	} else if (isUsernameAlreadyInUse) {
		return nil, errors.New("Username already in use")
	}

	// user
	var user entities.User

	// find user
	if err := s.userRepository.FindUserByIDWithAllHisProducts(claims.UserID, &user); err != nil {
		return nil, errors.New("Couldn't find user")
	}

	// delete related products from cache
	for i := range user.Products {
		s.redis.DEL(ctx, fmt.Sprintf("products:%d", user.Products[i].ID))
	}

	// update username
	user.Username = request.Username

	// update user in database
	if err := s.userRepository.Save(&user); err != nil {
		return nil, errors.New("Couldn't save updated user to database")
	}

	return &dtos.AuthenticationTokens{
		RefreshToken: tokens.RefreshToken,
		AccessToken: tokens.AccessToken,
	}, nil
}