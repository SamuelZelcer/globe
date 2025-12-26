package userService

import (
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/transaction"
	"globe/internal/service/jwt"
)

type Service interface {
	Create(request *dtos.SignUpRequest) (*string, error)
	Verify(request *dtos.VerifyUserRequest, token *string) error
}

type service struct {
	userRepository user.Repository
	unverifiedUserRepository unverifiedUser.Repository
	transactions transaction.Transactions
	jwtManager jwt.JWTManager
}

func InitService(
	userRepository user.Repository,
	unverifiedUserRepository unverifiedUser.Repository,
	transactions transaction.Transactions,
	jwtManager jwt.JWTManager,
) Service {
	return &service{
		userRepository: userRepository,
		unverifiedUserRepository: unverifiedUserRepository,
		transactions: transactions,
		jwtManager: jwtManager,
	}
}