package search

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
	"regexp"
)

var reg = regexp.MustCompile(`[\s'"]`)

type SearchService interface {
	Search(ctx context.Context, request *dtos.SearchRequest) (*dtos.SearchProductResponse, error)
}

type searchService struct {
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
	redis redis.Cache,
	jwtManager JWT.Manager,
	refreshTokenService refreshTokenService.Service,
) SearchService {
	return &searchService{
		productRepository: productRepository,
		userRepository: userRepository,
		email: email,
		redis: redis,
		jwtManager: jwtManager,
		refreshTokenService: refreshTokenService,
	}
}