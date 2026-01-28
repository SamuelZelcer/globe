package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"strings"
	"time"
)

func (s *searchService) Search(ctx context.Context, request *dtos.SearchRequest) (*dtos.SearchProductResponse, error) {
	// validate request
	if request.Name == "" ||
	request.Page == 0 {
		return nil, errors.New("Invalid request")
	}

	// parse product name
	parsedName := reg.ReplaceAllString(strings.ToLower(request.Name), "")
	
	// find a slice of products
	productIDsJSON, err := s.redis.GET(ctx, fmt.Sprintf("search:%s:page:%d", parsedName, request.Page))
	if err != nil {

		// find products from database
		return s.findProductsFromDB(
			ctx, request, parsedName,
		)
	}

	// reset TTL for productIDs
	if err := s.redis.EXPIRE(
		ctx,
		fmt.Sprintf("search:%s:page:%d", parsedName, request.Page),
		time.Hour*4,
	); err != nil {
		return nil, errors.New("Couldn't reset TTL for productIDs")
	}
	
	// parse product IDs from JSON
	var productIDs []uint64
	if err := json.Unmarshal([]byte(productIDsJSON), &productIDs); err != nil {
		return nil, errors.New("Couldn't parse IDs from JSON")
	}
	
	// parse IDs to redis keys
	var productsKeys []string
	for i := range productIDs {
		productsKeys = append(productsKeys, fmt.Sprintf("product:%d", productIDs[i]))
	}
	
	// get products
	productsRAW, err := s.redis.GETMORETHEN1(ctx, productsKeys)
	if err != nil {
		// find products from database
		return s.findProductsFromDB(
			ctx, request, parsedName,
		)
	}

	// parse raw products to entities
	var products []dtos.SearchProduct
	for i := range productsRAW {

		if productsRAW[i] == nil {
			product, err := s.findProductFromDB(ctx, productIDs[i])
			if err != nil {
				continue
			}
			products = append(products, dtos.SearchProduct{
				ProductID: product.ID,
				Name: product.Name,
				Price: fmt.Sprintf("%.2f", float64(product.Price)/100),
			})
			continue
		}

		// parse product from JSON to entity
		var parsedProduct entities.Product
		if err := json.Unmarshal([]byte(productsRAW[i].(string)), &parsedProduct); err != nil {
			product, err := s.findProductFromDB(ctx, productIDs[i])
			if err != nil {
				continue
			}
			products = append(products, dtos.SearchProduct{
				ProductID: product.ID,
				Name: product.Name,
				Price: fmt.Sprintf("%.2f", float64(product.Price)/100),
			})
			continue
		}

		// add product to slice
		products = append(
			products, dtos.SearchProduct{
				ProductID: productIDs[i],
				Name: parsedProduct.OriginalName,
				Price: fmt.Sprintf("%.2f", float64(parsedProduct.Price)/100),
			},
		)

		// reset TTL for product
		if err := s.redis.EXPIRE(
			ctx,
			fmt.Sprintf("product:%d", productIDs[i]),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't reset TTL for productIDs")
		}
	}
	// other essential things
	return s.findOtherStuff(ctx, parsedName, request, &products)
}