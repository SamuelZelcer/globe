package userService

import (
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
)

type Service interface {
	SignUp(request *dtos.SignUpRequest) (*string, error)
}

type service struct {
	userRepository user.Repository
	unverifiedUserRepository unverifiedUser.Repository
	email email.Email
	jwtManager JWT.Manager
}

func Init(
	userRepository user.Repository,
	unverifiedUserRepository unverifiedUser.Repository,
	email email.Email,
	jwtManager JWT.Manager,
) Service {
	return &service{
		userRepository: userRepository,
		unverifiedUserRepository: unverifiedUserRepository,
		email: email,
		jwtManager: jwtManager,
	}
}