package userService

import (
	"errors"
	"globe/internal/repository/dtos"
	"globe/internal/repository/entities"
	"math/rand"
	"net/mail"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(request *dtos.SignUpRequest) (*string, error) {
	// delete expired unverified users
	s.unverifiedUserRepository.DeleteExpiredUsers(time.Now())

	// validate request
	if request.Username == nil ||
	request.Email == nil ||
    request.Password == nil ||
	len(*request.Username) < 4 ||
	len(*request.Username) > 60 ||
	len(*request.Password) < 8 ||
	len(*request.Password) > 120 {
	    return nil, errors.New("Bad request")
	}

	// validate email
	if _, err := mail.ParseAddress(*request.Email); err != nil {
		return nil, errors.New("Bad request")
	}

	// validate username
	matched, err := regexp.MatchString(`[!@#$%^&*()]`,* request.Username)
	if err != nil || matched {
		return nil, errors.New("Invalid username")
	}

	// is user with username or email already exists
	if is := s.unverifiedUserRepository.IsUsernameOrEmailAlreadyInUse(request.Username, request.Email); is {
		return  nil, errors.New("User with this username or email already exists")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), 10)
	if err != nil {
		return nil, errors.New("Couldn't generate hashed passowd")
	}

	// unverified user
	user := &entities.UnverifiedUser{
		Username: *request.Username,
		Email: *request.Email,
		Password: hashedPassword,
		Expired: time.Now().Add(time.Minute * 15),
		Code: strconv.Itoa(rand.Intn(900000) + 100000),
	}

	// send email with verification code
	go s.email.SendVerificationCode(&user.Code, request.Email)

	// create unverified user and get his ID
	userID, err := s.unverifiedUserRepository.Create(user)
	if err != nil {
		return nil, errors.New("Couldn't create unverified user")
	}

	// generate jwt token with username and ID
	token, err := s.jwtManager.Create(userID, &user.Email, time.Minute*15)
	if err != nil {
		return nil, errors.New("Couldn't generate jwt token")
	}
	return token, nil
}