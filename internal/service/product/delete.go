package productService

import (
	"errors"
	"globe/internal/repository/dtos"
)

func (s *service) Delete(request *dtos.DeleteProductRequest, token *string) error {
	// validate request
	if request.ProductID == nil {
		return errors.New("Invalid request")
	}

	// validate jwt token and get claims
	_, err := s.jwtManager.Validate(token)
	if err != nil{
		return errors.New("Invalid token")
	}

	// delete product
	if err := s.productRepository.DeleteByID(request.ProductID); err != nil {
		return errors.New("Couldn't delete product")
	}
	return nil
}