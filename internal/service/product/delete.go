package productService

import (
	"context"
	"errors"
	"globe/internal/repository/dtos"
)

func (s *service) Delete(
	ctx context.Context,
	request *dtos.DeleteProductRequest,
	token *string,
) (*dtos.AuthenticationTokens, error) {
	// validate request
	if request.ProductID == nil {
		return nil, errors.New("Invalid request")
	}

	// validate token | update auth tokens
	tokens := &dtos.AuthenticationTokens{}
	_, err := s.jwtManager.Validate(token)
	if err != nil {
		if request.RefreshToken == nil {
			return nil, errors.New("Invalid jwt token")
		}
		_, err = s.refreshTokenService.Update(ctx, request.RefreshToken, token, tokens)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	// delete product
	if err := s.productRepository.DeleteByID(request.ProductID); err != nil {
		return nil, errors.New("Couldn't delete product")
	}
	return tokens, nil
}