package userService

import (
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"math/rand"
	"strconv"
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

func (s *service) GetNewCode(token *string) error {
	// get claims and validate token
	claims, err := s.jwtManager.VerifyAndGetClaims(token)
	if err != nil {
		return errors.New("Invalid token")
	}

	// update verification code
	newCode := strconv.Itoa(rand.Intn(900000) + 100000)
	if err := s.unverifiedUserRepository.UpdateVerificationCode(&claims.UserID, &newCode);
	err != nil {
		return errors.New("Couldn't update verificatio code")
	}

	// send new verification code to user
	go s.email.SendVerificationCode(&newCode, &claims.Subject)
	return nil
}

func (s *service) SendCodeAgain(token *string) error {
	// get claims and validate token
	claims, err := s.jwtManager.VerifyAndGetClaims(token)
	if err != nil {
		return errors.New("Invalid token")
	}

	// find usesr code
	var code string
	if err := s.unverifiedUserRepository.FindCodeByID(&claims.UserID, &code); err != nil {
		return errors.New("Couldn't find user's verification code")
	}

	// send verification ode again
	go s.email.SendVerificationCode(&code, &claims.Subject)
	return nil
}