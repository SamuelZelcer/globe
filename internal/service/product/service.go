package productService

import (
	"context"
	"globe/internal/repository/dtos"
	serviceDtos "globe/internal/repository/dtos/service"
	"globe/internal/repository/entities/product"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	"globe/internal/service/product/search"
	refreshTokenService "globe/internal/service/refreshToken"
	"regexp"
)


var (
	reg = regexp.MustCompile(`[\s'"]`)
	PRICEREGEXP = regexp.MustCompile(`^\d+\.\d{2}$`)
)

type Service interface {
	search.SearchService
	Create(
		ctx context.Context,
		request *dtos.CreateProductRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
	Update(
		ctx context.Context,
		request *dtos.UpdateProductRequest,
		token string,
	) (*dtos.AuthenticationTokens, *dtos.UpdateProductResponce, error)
	Delete(
		ctx context.Context,
		request *dtos.DeleteProductRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
}

type service struct {
	search.SearchService
	productRepository product.Repository
	userRepository user.Repository
	email email.Email
	redis redis.Cache
	jwtManager JWT.Manager
	refreshTokenService refreshTokenService.Service
}

func Init(d *serviceDtos.ProductDependencies) Service {
	return &service{
		SearchService: search.Init(d),
		productRepository: d.ProductRepository,
		userRepository: d.UserRepository,
		email: d.Email,
		redis: d.Redis,
		jwtManager: d.JWTManager,
		refreshTokenService: d.RefreshTokenService,
	}
}