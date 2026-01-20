package productService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
)

func (s *service) Delete(
	ctx context.Context,
	request *dtos.DeleteProductRequest,
	token *string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.ProductID == nil {
		return nil, errors.New("Invalid request")
	}

	// validate token
	tokens := &dtos.AuthenticationTokens{}
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == nil {
			return nil, errors.New("Invalid jwt token")
		}

		// update authentication tokens
		claims, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, tokens)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	// product
	product := &entities.Product{}
	
	// find product in redis
	productJSON, err := s.redis.GET(ctx, fmt.Sprintf("product:%d", *request.ProductID))
	if err != nil {

		// find product in database
		if err := s.productRepository.FindByID(request.ProductID, product); err != nil {
			return nil, errors.New("This product does not exist")
		}
	}

	// parse productJSON if find
	if productJSON != "" {
		if err := json.Unmarshal([]byte(productJSON), product); err != nil {
			return nil, errors.New("Couldn't parse product from JSON")
		}
	}

	// is this a user-owned product
	if claims.UserID != product.Owner {
		return nil, errors.New("It is not user owned product")
	}

	// delete product from redis
	s.redis.DEL(ctx, fmt.Sprintf("product:%d", *request.ProductID))

	// delete product from database
	if err := s.productRepository.DeleteByID(request.ProductID); err != nil {
		return nil, errors.New("Couldn't delete product")
	}
	
	return tokens, nil
}