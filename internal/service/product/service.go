package productService

import (
	"context"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/product"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	refreshTokenService "globe/internal/service/refreshToken"
)

type Service interface {
	Create(
		ctx context.Context,
		request *dtos.CreateProductRequest,
		token *string,
	) (*dtos.AuthenticationTokens, error)
	Update(
		ctx context.Context,
		request *dtos.UpdateProductRequest,
		token *string,
	) (*dtos.AuthenticationTokens, *dtos.UpdateProductResponse, error)
	Delete(
		ctx context.Context,
		request *dtos.DeleteProductRequest,
		token *string,
	) (*dtos.AuthenticationTokens, error)
}

type service struct {
	productRepository product.Repository
	userRepository user.Repository
	email email.Email
	transactions transactions.Transactions
	redis redis.Cache
	jwtManager JWT.Manager
	refreshTokenService refreshTokenService.Service
}

func Init(
	productRepository product.Repository,
	userRepository user.Repository,
	email email.Email,
	transactions transactions.Transactions,
	redis redis.Cache,
	jwtManager JWT.Manager,
	refreshTokenService refreshTokenService.Service,
) Service {
	return &service{
		productRepository: productRepository,
		userRepository: userRepository,
		email: email,
		transactions: transactions,
		redis: redis,
		jwtManager: jwtManager,
		refreshTokenService: refreshTokenService,
	}
}