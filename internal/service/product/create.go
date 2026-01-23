package productService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math"
	"strconv"
	"strings"
	"time"
)

func (s *service) Create(
	ctx context.Context,
	request *dtos.CreateProductRequest,
	token string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Name == "" ||
	request.Price == "" ||
	request.Description == "" ||
	len(request.Name) > 100 ||
	len(request.Description) > 800 {
		return nil, errors.New("Invalid request")
	}

	// validate token and get claims
	var tokens dtos.AuthenticationTokens
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == "" {
			return nil, errors.New("Invalid jwt token")
		}

		// update authentication tokens
		claims, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, &tokens)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	// parse product name
	parsedName := reg.ReplaceAllString(strings.ToLower(request.Name), "")
	
	// validate price .00
	if !PRICEREGEXP.MatchString(request.Price) {
		return nil, errors.New("Invalid price")
	}
	
	// convert price to float value
	floatPrice, err := strconv.ParseFloat(request.Price, 64)
	if err != nil {
		return nil, errors.New("Couldn't convert price to float value")
	}

	// find user
	var user entities.User
	if err := s.userRepository.FindByiD(claims.UserID, &user); err != nil {
		return nil, errors.New("couldn't find user")
	}
	
	// product
	product := entities.Product{
		Name: parsedName,
		OriginalName: request.Name,
		Price: uint64(math.Round(floatPrice*100)),
		Description: request.Description,
		Owner: claims.UserID,
		User: user,
	}

	// save product to database
	productID, err := s.productRepository.Save(&product)
	if err != nil {
		return nil, errors.New("Couldn't save product")
	}
	
	// parse product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, errors.New("Couldn't parse product to JSON")
	}
	
	// save product to redis
	if err := s.redis.SET(
		ctx,
		fmt.Sprintf("product:%d", productID),
		string(productJSON),
		time.Hour*4,
	); err != nil {
		return nil, errors.New("Couldn't save product to redis")
	}

	// delete most popular pages from cache
	for i := range 5 {
		s.redis.DEL(ctx, fmt.Sprintf("search:%s:page:%d", parsedName, i))
	}

	return &tokens, nil
}