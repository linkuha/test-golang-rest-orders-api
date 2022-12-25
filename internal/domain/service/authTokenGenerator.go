package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	signingKey = "asdq2ed1212d12dawsx"
	tokenTTL   = 12 * time.Hour
)

type AuthTokenGenerator struct {
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func (t AuthTokenGenerator) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	return token.SignedString([]byte(signingKey))
}

func (t AuthTokenGenerator) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}
