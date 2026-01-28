package search

import (
	"context"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"strconv"
	"time"
)

func (s *searchService) findOtherStuff(
	ctx context.Context,
	parsedName string,
	request *dtos.SearchRequest,
	products *[]dtos.SearchProduct,
) (*dtos.SearchProductResponse, error) {
	// get amount of products from redis
	var amountOfProducts int64
	amountJSON, err := s.redis.GET(
		ctx,
		fmt.Sprintf("search:%s", parsedName),
	)
	if err != nil {
		// get amount of products from database
		if err := s.productRepository.CountProducts(parsedName, &amountOfProducts); err != nil {
			return nil, errors.New("Couldn't count products")
		}

		if amountOfProducts >= 1 {
			// save amount of products to redis
			if err := s.redis.SET(
				ctx,
				fmt.Sprintf("search:%s", parsedName),
				fmt.Sprintf("%d", amountOfProducts),
				time.Hour*4,
			); err != nil {
				return nil, errors.New("Couldn't save amount of products to redis")
			}
		}
	}

	// reset TTL for amountOfProducts
	if err := s.redis.EXPIRE(
		ctx,
		fmt.Sprintf("search:%s", parsedName),
		time.Hour*4,
	); err != nil {
		return nil, errors.New("Couldn't reset TTL for productIDs")
	}

	// parse amountJSON if find
	if amountJSON != "" {
		amountOfProducts, err = strconv.ParseInt(amountJSON, 10, 64)
		if err != nil {
			return nil, errors.New("Couldn't parse amountJSON")
		}
	}

	return &dtos.SearchProductResponse{
		TotalAmountOfProducts: amountOfProducts,
		TotalAmountOfPages: (amountOfProducts+14)/15,
		CurrentPage: request.Page,
		Products: products,
	}, nil
}