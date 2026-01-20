package productService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"strconv"
	"strings"
	"time"
)

func (s *service) Search(ctx context.Context, request dtos.SearchRequest) (*dtos.SearchProductResponse, error) {
	// validate request
	if request.Name == nil ||
	request.Page == nil {
		return nil, errors.New("Invalid request")
	}

	// parse product name
	parsedName := reg.ReplaceAllString(strings.ToLower(*request.Name), "")
	
	// find a slice of products
	productIDsJSON, err := s.redis.GET(ctx, fmt.Sprintf("search:%s:page:%d", parsedName, *request.Page))
	if err != nil {

		// find products from database
		return s.findProductsFromDB(
			ctx, request, &parsedName,
		)
	}

	// reset TTL for productIDs
	if err := s.redis.EXPIRE(
		ctx,
		fmt.Sprintf("search:%s:page:%d", parsedName, *request.Page),
		time.Hour*4,
	); err != nil {
		return nil, errors.New("Couldn't reset TTL for productIDs")
	}
	
	// parse product IDs from JSON
	var productIDs []uint64
	if err := json.Unmarshal([]byte(productIDsJSON), &productIDs); err != nil {
		return nil, errors.New("Couldn't parse IDs from JSON")
	}
	
	// products
	var products []dtos.SearchProduct

	// find each product
	for _, v := range productIDs {
		product := &entities.Product{}

		// find product in cache
		productJSON, err := s.redis.GET(ctx, fmt.Sprintf("product:%d", v))
		if err != nil {

			// find product in database
			if err := s.productRepository.FindByID(&v, product); err != nil {
				
				// delete outdated page from redis
				if err := s.redis.DEL(
					ctx,
					fmt.Sprintf("search:%s:page:%d", parsedName, *request.Page),
				); err != nil {
					return nil, errors.New("Couldn't delete outdated page from redis")
				}

				// find products from database
				return s.findProductsFromDB(
					ctx, request, &parsedName,
				)
			}
			
			// parse product to JSON
			productJSON, err := json.Marshal(product)
			if err != nil {
				return nil, errors.New("Couldn't parse product to JSON")
			}
			
			if string(productJSON) != "" && string(productJSON) != "null" {
				// save product to redis
				if err := s.redis.SET(
					ctx,
					fmt.Sprintf("product:%d", v),
					string(productJSON),
					time.Hour*4,
				); err != nil {
					return nil, errors.New("Couldn't save product to redis")
				}
			}
		}

		// reset TTL for product
		if err := s.redis.EXPIRE(
			ctx,
			fmt.Sprintf("product:%d", v),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't reset TTL for productIDs")
		}

		// parse productJSON if find
		if productJSON != "" && productJSON != "null" {
			if err := json.Unmarshal([]byte(productJSON), product); err != nil {
				return nil, errors.New("Couldn't parse product from JSON")
			}
		}

		// convert price
		displayPrice := fmt.Sprintf("%.2f", float64(product.Price)/100)
			
		// add product to slice
		products = append(
			products,
			dtos.SearchProduct{
				ProductID: &product.ID,
				Name: &product.Name,
				Price: &displayPrice,
			},
		)
	}

	// get amount of products from redis
	var amountOfProducts int64
	amountJSON, err := s.redis.GET(
		ctx,
		fmt.Sprintf("search:%s", parsedName),
	)
	if err != nil {
		// get amount of products from database
		if err := s.productRepository.CountProducts(&parsedName, &amountOfProducts); err != nil {
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
	amountOfPages := (amountOfProducts+14)/15

	return &dtos.SearchProductResponse{
		TotalAmountOfProducts: &amountOfProducts,
		TotalAmountOfPages: &amountOfPages,
		CurrentPage: request.Page,
		Products: &products,
	}, nil
}

func (s *service) findProductsFromDB(
	ctx context.Context,
	request dtos.SearchRequest,
	parsedName *string,
) (*dtos.SearchProductResponse, error) {
	// product entities
	var productEntities []entities.Product

	// products
	var products []dtos.SearchProduct
	
	// find products from database
	if err := s.productRepository.FindProductsForSearch(
		parsedName,
		(int(*request.Page)-1)*15,
		&productEntities,
	); err != nil {
		return nil, errors.New("Couldn't find products")
	}

	var productIDs []uint64
	for _, v := range productEntities {
		// append product ID to IDs slice
		productIDs = append(productIDs, v.ID)

		// convert price
		displayPrice := fmt.Sprintf("%.2f", float64(v.Price)/100)
			
		// mapped productEmtity to prodicts slice
		products = append(products, dtos.SearchProduct{
			ProductID: &v.ID,
			Name: &v.Name,
			Price: &displayPrice,
		})

		// parse product to JSON
		productJSON, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("Couldn't parse product to JSON")
		}
			
		if string(productJSON) != "" && string(productJSON) != "null" {
			// save product to redis
			if err := s.redis.SET(
				ctx,
				fmt.Sprintf("product:%d", v.ID),
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
			fmt.Sprintf("search:%s:page:%d", *parsedName, *request.Page),
			string(productIDsJSON),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't save product to redis")
		}
	}

	// count products by provided name
	var amountOfProducts int64
	if err := s.productRepository.CountProducts(parsedName, &amountOfProducts); err != nil {
		return nil, errors.New("Couldn't count products")
	}
	amountOfPages := (amountOfProducts+14)/15

	if amountOfProducts >= 1 {
		// save amount of products to redis
		if err := s.redis.SET(
			ctx,
			fmt.Sprintf("search:%s", *parsedName),
			fmt.Sprintf("%d", amountOfProducts),
			time.Hour*4,
		); err != nil {
			return nil, errors.New("Couldn't save amount of products to redis")
		}
	}

	return &dtos.SearchProductResponse{
		TotalAmountOfProducts: &amountOfProducts,
		TotalAmountOfPages: &amountOfPages,
		CurrentPage: request.Page,
		Products: &products,
	}, nil
}