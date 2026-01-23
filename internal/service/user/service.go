package userService

import (
	"context"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/refreshToken"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
)

type Service interface {
	SignUp(request *dtos.SignUpRequest) (string, error)
	Verification(request *dtos.VerifyUserRequest, token string) error
	GetNewCode(token string) error
	SendCodeAgain(token string) error
	SignIn(request *dtos.SignInRequest, ctx context.Context) (*dtos.AuthenticationTokens, error)
}

type service struct {
	userRepository user.Repository
	unverifiedUserRepository unverifiedUser.Repository
	email email.Email
	jwtManager JWT.Manager
	transactions transactions.Transactions
	redis redis.Cache
	refreshTokenRepository refreshToken.Repository
}

func Init(
	userRepository user.Repository,
	unverifiedUserRepository unverifiedUser.Repository,
	email email.Email,
	jwtManager JWT.Manager,
	transactions transactions.Transactions,
	redis redis.Cache,
	refreshTokenRepository refreshToken.Repository,
) Service {
	return &service{
		userRepository: userRepository,
		unverifiedUserRepository: unverifiedUserRepository,
		email: email,
		jwtManager: jwtManager,
		transactions: transactions,
		redis: redis,
		refreshTokenRepository: refreshTokenRepository,
	}
}