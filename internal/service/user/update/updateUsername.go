package userUpdateService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"regexp"
)

func (s *userUpdateService) UpdateUsername(
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

	// find user
	var user entities.User
	if err := s.userRepository.FindUserByIDWithAllHisProducts(claims.UserID, &user); err != nil {
		return nil, errors.New("Couldn't find user")
	}

	// delete related products from cache
	var keys []string
	for i := range user.Products {
		keys = append(keys, fmt.Sprintf("product:%d", user.Products[i].ID))
	}
	if err := s.redis.DELMORETHEN1(ctx, keys); err != nil {
		return nil, errors.New("Couldn't delete outdated data from redis")
	}

	// update username
	user.Username = request.Username

	// update user in database
	if err := s.userRepository.Save(&user); err != nil {
		return nil, errors.New("Couldn't save updated user to database")
	}
	return &tokens, nil
}