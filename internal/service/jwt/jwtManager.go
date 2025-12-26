package jwt

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	Create(userID *uint32, username *string, duration time.Duration) (*string, error)
	VerifyAndGetClaims(tokenStr *string) (*UserClaims, error)
}

type jwtManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
}

func InitJWTManager() JWTManager {
	privateKey, err := LoadOrGenerateECDSAKeys()
	if err != nil {
		log.Fatalf("Couldn't load or generate ecdsa keys %v\n", err)
	}
	return &jwtManager{
		privateKey: privateKey,
		publicKey: &privateKey.PublicKey,
	}
}

func (m *jwtManager) Create(userID *uint32, username *string, duration time.Duration) (*string, error) {
	claims, err := IntiUserClaims(userID, username, duration)
	if err != nil {
		return nil, err
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token, err := unsignedToken.SignedString(m.privateKey)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (m *jwtManager) VerifyAndGetClaims(tokenStr *string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(*tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, errors.New("Invalid token signing method")
		}
		return m.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return claims, errors.New("Invalid token claims")
	}
	return claims, nil
}