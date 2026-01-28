package search

import (
	"context"
	"globe/internal/repository/dtos"
	serviceDTOs "globe/internal/repository/dtos/service"
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

func Init(d *serviceDTOs.ProductDependencies) SearchService {
	return &searchService{
		productRepository: d.ProductRepository,
		userRepository: d.UserRepository,
		email: d.Email,
		redis: d.Redis,
		jwtManager: d.JWTManager,
		refreshTokenService: d.RefreshTokenService,
	}
}