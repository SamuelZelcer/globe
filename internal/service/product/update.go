package productService

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math"
	"strconv"
)

func (s *service) Update(
	ctx context.Context,
	request *dtos.UpdateProductRequest,
	token *string,
) (*dtos.AuthenticationTokens, *dtos.UpdateProductResponse, error) {
	// validate request
	if request.ProductID == nil {
		return nil, nil, errors.New("Invalid request")
	}

	// validate token | update auth tokens
	tokens := &dtos.AuthenticationTokens{}
	_, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == nil {
			return nil, nil, errors.New("Invalid jwt token")
		}
		_, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, tokens)
		if err != nil {
			return nil, nil, errors.New(err.Error())
		}
	}

	// find product
	var product entities.Product
	if err := s.productRepository.FindByID(request.ProductID, &product); err != nil{
		return nil, nil, errors.New("Couldn't find product")
	}

	// update name
	if request.Name != nil {
		product.Name = *request.Name
	}

	// update price
	if request.Price != nil {
		// validate price .00
		if !PRICEREGEXP.MatchString(*request.Price) {
			return nil, nil, errors.New("Invalid price")
		}

		// convert price to float value
		floatPrice, err := strconv.ParseFloat(*request.Price, 32)
		if err != nil {
			return nil, nil, errors.New("Couldn't convert price to float value")
		}
		product.Price = uint64(math.Round(floatPrice*100))
	}

	// update description
	if request.Description != nil {
		product.Description = *request.Description
	}

	// save updated product
	if err := s.productRepository.Save(&product); err != nil {
		return nil, nil, errors.New("Couldn't save updated product")
	}
	
	// convert prive
	displayPrice := fmt.Sprintf("%.2f", float64(product.Price)/100)

	return tokens, 
	&dtos.UpdateProductResponse{
		Name: &product.Name,
		Price: &displayPrice,
		Description: &product.Description,
	}, nil
}