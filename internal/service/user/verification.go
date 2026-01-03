package userService

import (
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
)

func (s *service) Verification(request *dtos.VerifyUserRequest, token *string) error {
	// validate request
	if request.Code == "" ||
	len(request.Code) != 6 {
		return errors.New("Bad request")
	}

	// verify token and get user claims
	claims, err := s.jwtManager.VerifyAndGetClaims(token)
	if err != nil {
		return err
	}

	// find unverified user bi ID
	unverifiedUser :=  &unverifiedUser.UnverifiedUser{}
	if err := s.unverifiedUserRepository.FindByID(&claims.UserID, unverifiedUser); err != nil {
		return err
	}

	// verify user's code
	if request.Code != unverifiedUser.Code {
		return errors.New("Incorrect verification code")
	}

	// user
	user := user.User{
		Username: unverifiedUser.Username,
		Email: unverifiedUser.Email,
		Password: unverifiedUser.Password,
	}

	// begin transaction
	tx := s.transactions.BeginTransaction()

	// make sure that transaction will rollback if something fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// delete unverified user
	if err := s.unverifiedUserRepository.DeleteByID(&unverifiedUser.ID); err != nil {
		tx.Rollback()
		return err
	}

	// create user
	if err := s.userRepository.Create(&user); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	if err := s.transactions.CommitTransaction(tx); err != nil {
		return err
	}
	return nil
} 