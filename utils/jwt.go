package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

type TokenClaims struct {
	Email  string
	UserID string
}

func GenerateToken(email, userID string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func VerifyToken(tokenString string) (*TokenClaims, error, *TokenClaims) {
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token"), nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse token claims"), nil
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token"), nil
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id not found in token"), nil
	}

	return &TokenClaims{UserID: userID, Email: email}, nil, nil
}
