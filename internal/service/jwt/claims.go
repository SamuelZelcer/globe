package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uint32
	jwt.RegisteredClaims
}

func IntiUserClaims(userID *uint32, username *string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	if username == nil {
		return &UserClaims{
			UserID: *userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ID: tokenID.String(),
				Subject: "",
				IssuedAt: jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			},
		}, nil
	}
	return &UserClaims{
		UserID: *userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenID.String(),
			Subject: *username,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}