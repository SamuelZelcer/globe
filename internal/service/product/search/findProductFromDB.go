package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"time"
)

func (s *searchService) findProductFromDB(
	ctx context.Context,
	productID uint64,
) (*entities.Product, error) {
	var product entities.Product

	// find product in database
	if err := s.productRepository.FindByID(productID, &product); err != nil {
		return nil, errors.New("Couldn't find product in database")
	}
	// parse product to JSON
	productJSON, err := json.Marshal(dtos.CachedProduct{
		Name: product.Name,
		OriginalName: product.OriginalName,
		Price: product.Price,
		Description: product.Description,
		Owner: product.User.ID,
		OwnerName: product.User.Username,
	})
	if err != nil {
		return nil, errors.New("Couldn't parse product to JSON")
	}

	if string(productJSON) != "" && string(productJSON) != "null" {
		// save product to redis
		if err := s.redis.SET(
			ctx,
			fmt.Sprintf("product:%d", product.ID),
			string(productJSON),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't save product to redis")
		}
	}
	return &product, nil
}