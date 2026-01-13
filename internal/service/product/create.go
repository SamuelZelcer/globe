package productService

import (
	"context"
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math"
	"regexp"
	"strconv"
)

var PRICEREGEXP = regexp.MustCompile(`^\d+\.\d{2}$`)

func (s *service) Create(
	ctx context.Context,
	request *dtos.CreateProductRequest,
	token *string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.Name == nil ||
	request.Price == nil ||
	request.Description == nil ||
	len(*request.Name) > 100 ||
	len(*request.Description) > 800 {
		return nil, errors.New("Invalid request")
	}

	// validate token and get claims | update auth tokens
	tokens := &dtos.AuthenticationTokens{}
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == nil {
			return nil, errors.New("Invalid jwt token")
		}
		claims, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, tokens)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	// validate price .00
	if !PRICEREGEXP.MatchString(*request.Price) {
		return nil, errors.New("Invalid price")
	}

	// convert price to float value
	floatPrice, err := strconv.ParseFloat(*request.Price, 64)
	if err != nil {
		return nil, errors.New("Couldn't convert price to float value")
	}

	// product
	product := &entities.Product{
		Name: *request.Name,
		Price: uint64(math.Round(floatPrice*100)),
		Description: *request.Description,
		Owner: claims.UserID,
	}

	// save product
	if err := s.productRepository.Save(product); err != nil {
		return nil, errors.New("Couldn't save product")
	}
	return tokens, nil
}