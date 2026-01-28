package serviceDTOs

import (
	"globe/internal/repository/entities/product"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	refreshTokenService "globe/internal/service/refreshToken"
)

type ProductDependencies struct {
	ProductRepository product.Repository
	UserRepository user.Repository
	Email email.Email
	Redis redis.Cache
	JWTManager JWT.Manager
	RefreshTokenService refreshTokenService.Service
}

type UserDependencies struct {
	UserRepository user.Repository
	UnverifiedUserRepository unverifiedUser.Repository
	Email email.Email
	JWTManager JWT.Manager
	Transactions transactions.Transactions
	Redis redis.Cache
	RefreshTokenService refreshTokenService.Service
}