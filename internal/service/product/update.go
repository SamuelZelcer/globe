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

func (s *service) Update(
	ctx context.Context,
	request *dtos.UpdateProductRequest,
	token *string,
) (*dtos.AuthenticationTokens, *dtos.UpdateProductResponce, error) {
	// validate request
	if request.ProductID == nil {
		return nil, nil, errors.New("Invalid request")
	}

	// validate token
	tokens := &dtos.AuthenticationTokens{}
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == nil {
			return nil, nil, errors.New("Invalid jwt token")
		}

		// update authentication tokens
		claims, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, tokens)
		if err != nil {
			return nil, nil, errors.New(err.Error())
		}
	}

	// product
	product := &entities.Product{}
	
	// find product in redis
	productJSON, err := s.redis.GET(ctx, fmt.Sprintf("product:%d", *request.ProductID))
	if err != nil {

		// find product in database
		if err := s.productRepository.FindByID(request.ProductID, product); err != nil {
			return nil, nil, errors.New("This product does not exist")
		}
		
	}

	// parse productJSON if find
	if productJSON != "" {
		if err := json.Unmarshal([]byte(productJSON), product); err != nil {
			return nil, nil, errors.New("Couldn't parse product from JSON")
		}
	}

	// parse product name
	var parsedName string
	if request.Name != nil {
		parsedName = reg.ReplaceAllString(strings.ToLower(*request.Name), "")
	}

	// is this a user-owned product
	if claims.UserID != product.Owner {
		return nil, nil, errors.New("It is not user owned product")
	}

	// update name if provided
	if request.Name != nil {
		
		// delete most popular pages from cache with old name
		for i := range 5 {
			s.redis.DEL(ctx, fmt.Sprintf("search:%s:page:%d", product.Name, i))
			s.redis.DEL(ctx, fmt.Sprintf("search:%s:page:%d", parsedName, i))
		}
		product.Name = parsedName
	}

	// update price if provided
	if request.Price != nil {
		// validate price .00
		if !PRICEREGEXP.MatchString(*request.Price) {
			return nil, nil, errors.New("Invalid price")
		}

		// convert price to float value
		floatPrice, err := strconv.ParseFloat(*request.Price, 64)
		if err != nil {
			return nil, nil, errors.New("Couldn't convert price to float value")
		}
		product.Price = uint64(math.Round(floatPrice*100))
	}

	// update description if provided
	if request.Description != nil {
		product.Description = *request.Description
	}

	// parse product to JSON
	newProductJSON, err := json.Marshal(product)
	if err != nil {
		return nil, nil,  errors.New("Couldn't parse updated product to JSON")
	}

	// save updated product to database
	if _, err := s.productRepository.Save(product); err != nil {
		return nil, nil, errors.New("Couldn't save updated product")
	}
	
	// save updated product to redis
	if err := s.redis.SET(
		ctx,
		fmt.Sprintf("product:%d", product.ID),
		string(newProductJSON),
		time.Hour*4,
	); err != nil {
		return nil, nil,  errors.New("Couldn't save updated product to redis")
	}

	// convert price
	displayPrice := fmt.Sprintf("%.2f", float64(product.Price)/100)

	return tokens,
	&dtos.UpdateProductResponce{
		Name: &product.Name,
		Price: &displayPrice,
		Description: &product.Description,
	}, nil
}