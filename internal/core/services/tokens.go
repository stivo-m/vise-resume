package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s TokenService) CreateToken(id string, expiryDate time.Time) (string, error) {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("unable to get the token secret from env variables")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": expiryDate.Unix(),
		},
	)

	key := []byte(secretKey)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s TokenService) VerifyToken(tokenString string) (string, error) {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("unable to get the token secret from env variables")
	}

	key := []byte(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("either token is invalid or has expired")
	}

	// Extract the claims from the token if needed
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := fmt.Sprintf("%v", claims["id"])
		return userId, nil
	}

	return "", errors.New("unknown error occurred")
}
