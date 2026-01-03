package JWT

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager interface {
	Create(userID *uint32, username *string, duration time.Duration) (*string, error)
	VerifyAndGetClaims(tokenStr *string) (*UserClaims, error)
}

type manager struct {
	privateKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
}

func Init() Manager {
	privateKey, err := LoadOrGenerateECDSAKeys()
	if err != nil {
		log.Fatalf("Couldn't load or generate ecdsa keys %v\n", err)
	}
	return &manager{
		privateKey: privateKey,
		publicKey: &privateKey.PublicKey,
	}
}

func (m *manager) Create(userID *uint32, email *string, duration time.Duration) (*string, error) {
	unsignedToken := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		InitUserClaims(userID, email, &duration),
	)
	token, err := unsignedToken.SignedString(m.privateKey)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (m *manager) VerifyAndGetClaims(tokenStr *string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(*tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
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