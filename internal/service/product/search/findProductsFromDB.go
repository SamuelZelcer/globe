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

func (s *searchService) findProductsFromDB(
	ctx context.Context,
	request *dtos.SearchRequest,
	parsedName string,
) (*dtos.SearchProductResponse, error) {
	// products
	var products []dtos.SearchProduct
	
	// product entities
	var productEntities []entities.Product

	// product IDs
	var productIDs []uint64

	// find products from database
	if err := s.productRepository.FindProductsForSearch(
		parsedName,
		(int(request.Page)-1)*15,
		&productEntities,
	); err != nil {
		return nil, errors.New("Couldn't find products")
	}

	for i := range productEntities {
		// append product ID to IDs slice
		productIDs = append(productIDs, productEntities[i].ID)

		// mapped productEmtity to prodicts slice
		products = append(products, dtos.SearchProduct{
			ProductID: productEntities[i].ID,
			Name: productEntities[i].OriginalName,
			Price: fmt.Sprintf("%.2f", float64(productEntities[i].Price)/100),
		})

		// parse product to JSON
		productJSON, err := json.Marshal(dtos.CachedProduct{
			Name: productEntities[i].Name,
			OriginalName: productEntities[i].OriginalName,
			Price: productEntities[i].Price,
			Description: productEntities[i].Description,
			Owner: productEntities[i].User.ID,
			OwnerName: productEntities[i].User.Username,
		})
		if err != nil {
			return nil, errors.New("Couldn't parse product to JSON")
		}
			
		if string(productJSON) != "" && string(productJSON) != "null" {
			// save product to redis
			if err := s.redis.SET(
				ctx,
				fmt.Sprintf("product:%d", productEntities[i].ID),
				string(productJSON),
				time.Hour*4,
			); err != nil {
				return nil, errors.New("Couldn't save product to redis")
			}
		}
	}
		
	if len(productIDs) >= 1 {
		// parse product IDs to JSON
		productIDsJSON, err := json.Marshal(productIDs)
		if err != nil {
			return nil, errors.New("Couldn't parse product to JSON")
		}

		// save page to redis
		if err := s.redis.SET(
			ctx,
			fmt.Sprintf("search:%s:page:%d", parsedName, request.Page),
			string(productIDsJSON),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't save product to redis")
		}
	}

	// count products
	var amountOfProducts int64
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

	return &dtos.SearchProductResponse{
		TotalAmountOfProducts: amountOfProducts,
		TotalAmountOfPages: (amountOfProducts+14)/15,
		CurrentPage: request.Page,
		Products: &products,
	}, nil
}