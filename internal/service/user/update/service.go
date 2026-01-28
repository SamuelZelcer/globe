package userUpdateService

import (
	"context"
	"globe/internal/repository/dtos"
	serviceDTOs "globe/internal/repository/dtos/service"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	refreshTokenService "globe/internal/service/refreshToken"
)

type UserUpdateService interface {
	UpdateUsername(
		ctx context.Context,
		request *dtos.UpdateUsernameRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
	UpdateEmail(
		ctx context.Context,
		request *dtos.UpdateEmailRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
	NewEmailVerification(
		ctx context.Context,
		request *dtos.VerificationRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
	UpdatePassword(
		ctx context.Context,
		request *dtos.UpdatePasswordRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
	NewPasswordVerification(
		ctx context.Context,
		request *dtos.VerificationRequest,
		token string,
	) (*dtos.AuthenticationTokens, error)
}

type userUpdateService struct {
	userRepository user.Repository
	unverifiedUserRepository unverifiedUser.Repository
	email email.Email
	jwtManager JWT.Manager
	transactions transactions.Transactions
	redis redis.Cache
	refreshTokenService refreshTokenService.Service
}

func Init(d *serviceDTOs.UserDependencies) UserUpdateService {
	return &userUpdateService{
		userRepository: d.UserRepository,
		unverifiedUserRepository: d.UnverifiedUserRepository,
		email: d.Email,
		jwtManager: d.JWTManager,
		transactions: d.Transactions,
		redis: d.Redis,
		refreshTokenService: d.RefreshTokenService,
	}
}