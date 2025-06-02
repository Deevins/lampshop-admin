package infra

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("my_super_secret_key")

const tokenTTL = time.Hour * 2

// GenerateJWT генерирует JWT-токен со сроком tokenTTL.
func GenerateJWT(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT проверяет токен, возвращает username из поля "sub" или ошибку.
func ParseJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", jwt.ErrTokenInvalidClaims
}
