package productService

import (
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math"
	"regexp"
	"strconv"
)

var PRICEREGEXP = regexp.MustCompile(`^\d+\.\d{2}$`)

func (s *service) Create(request *dtos.CreateProductRequest, token *string) error {
	// validate request
	if request.Name == nil ||
	request.Price == nil ||
	request.Description == nil ||
	len(*request.Name) > 100 ||
	len(*request.Description) > 800 {
		return errors.New("Invalid request")
	}

	// validate token and get claims
	claims, err := s.jwtManager.Validate(token)
	if err != nil {
		return errors.New("invalid jwt token")
	}

	// validate price .00
	if !PRICEREGEXP.MatchString(*request.Price) {
		return errors.New("Invalid price")
	}

	// convert price to float value
	floatPrice, err := strconv.ParseFloat(*request.Price, 64)
	if err != nil {
		return errors.New("Couldn't convert price to float value")
	}

	// product
	product := &entities.Product{
		Name: *request.Name,
		Price: uint32(math.Round(floatPrice*100)),
		Description: *request.Description,
		Owner: claims.UserID,
	}

	// save product
	if err := s.productRepository.Save(product); err != nil {
		return errors.New("Couldn't save product")
	}
	return nil
}