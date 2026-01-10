package productService

import (
	"errors"
	"fmt"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math"
	"strconv"
)

func (s *service) Update(request *dtos.UpdateProductRequest, token *string) (*dtos.UpdateProductResponse, error) {
	// validate request
	if request.ProductID == nil {
		return nil, errors.New("Invalid request")
	}

	// validate token and get claims
	_, err := s.jwtManager.Validate(token)
	if err != nil {
		return nil, errors.New("Invalid token")
	}

	// find product
	var product entities.Product
	if err := s.productRepository.FindByID(request.ProductID, &product); err != nil{
		return nil, errors.New("Couldn't find product")
	}

	// update name
	if request.Name != nil {
		product.Name = *request.Name
	}

	// update price
	if request.Price != nil {
		// validate price .00
		if !PRICEREGEXP.MatchString(*request.Price) {
			return nil, errors.New("Invalid price")
		}

		// convert price to float value
		floatPrice, err := strconv.ParseFloat(*request.Price, 32)
		if err != nil {
			return nil, errors.New("Couldn't convert price to float value")
		}
		product.Price = uint32(math.Round(floatPrice*100))
	}

	// update description
	if request.Description != nil {
		product.Description = *request.Description
	}

	// save updated product
	if err := s.productRepository.Save(&product); err != nil {
		return nil, errors.New("Couldn't save updated product")
	}
	
	// convert prive
	displayPrice := fmt.Sprintf("%.2f", float64(product.Price)/100)

	return &dtos.UpdateProductResponse{
		Name: &product.Name,
		Price: &displayPrice,
		Description: &product.Description,
	}, nil
}