package userService

import (
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/service/email"
	"math/rand"
	"net/mail"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) Create(request *dtos.SignUpRequest) (*string, error) {
	// delete expired unverified users
	s.unverifiedUserRepository.DeleteExpired(time.Now())
	
	// validate request
	if request.Username == "" ||
	   request.Email == "" ||
	   request.Password == "" ||
	   len(request.Username) < 4 ||
	   len(request.Username) > 60 ||
	   len(request.Password) < 8 ||
	   len(request.Password) > 120 {
	    return nil, errors.New("Bad request")
	}
	
	// validate email
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, errors.New("Bad request")
	}

	// validate username
	matched, err := regexp.MatchString(`[!@#$%^&*()]`, request.Username)
	if err != nil {
		return nil, errors.New("Couldn't validate username")
	} else if matched {
		return nil, errors.New("Invalid username")
	}

	// check if username or email already in use
	usernameOrEmailExists, usernameErr, emailErr := s.userRepository.IsUsernameOrEmailAlreadyInUse(
		&request.Username, &request.Email,
	)
	if (usernameErr != gorm.ErrRecordNotFound || emailErr != gorm.ErrRecordNotFound) && 
	(usernameErr != nil || emailErr != nil) {
		return nil, errors.New("Couldn't check username and email for existing")
	} else if usernameOrEmailExists {
		return nil, errors.New("User with provided username or email already exists")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, errors.New("Couldn't generate hashed passowd")
	}

	// unverified user
	user := unverifiedUser.UnverifiedUser {
		Username: request.Username,
		Email: request.Email,
		Password: hashedPassword,
		Expired: time.Now().Add(time.Minute * 20),
		Code: strconv.Itoa(rand.Intn(900000) + 100000),
	}

	// send email with verification code
	go email.Send(&request.Email, "verification code", &user.Code)

	// create unverified user and get his ID
	userID, err := s.unverifiedUserRepository.Create(&user)
	if err != nil {
		return nil, errors.New("Couldn't create unverified user")
	}

	// generate jwt token with username and ID
	token, err := s.jwtManager.Create(&userID, nil, time.Minute*15)
	if err != nil {
		return nil, errors.New("Couldn't generate jwt token")
	}
	return token, nil
}