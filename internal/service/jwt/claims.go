package JWT

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uint64
	jwt.RegisteredClaims
}

func InitUserClaims(userID uint64, email string, duration *time.Duration) *UserClaims {
	return &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.NewString(),
			Subject: email,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(*duration)),
		},
	}
}