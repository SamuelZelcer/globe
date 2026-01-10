package productService

import (
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/product"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
)

type Service interface {
	Create(request *dtos.CreateProductRequest, token *string) error
	Update(request *dtos.UpdateProductRequest, token *string) (*dtos.UpdateProductResponse, error)
}

type service struct {
	productRepository product.Repository
	userRepository user.Repository
	email email.Email
	transactions transactions.Transactions
	redis redis.Repository
	jwtManager JWT.Manager
}

func Init(
	productRepository product.Repository,
	userRepository user.Repository,
	email email.Email,
	transactions transactions.Transactions,
	redis redis.Repository,
	jwtManager JWT.Manager,
) Service {
	return &service{
		productRepository: productRepository,
		userRepository: userRepository,
		email: email,
		transactions: transactions,
		redis: redis,
		jwtManager: jwtManager,
	}
}