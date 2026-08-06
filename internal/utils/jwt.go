package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func InitJWT(key string) {
	jwtKey = []byte(key)
}

type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, email, role string) (string, error) {
	if len(jwtKey) == 0 {
		return "", errors.New("jwt key not initialized")
	}

	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprint(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenStr string) (int64, string, string, error) {
	if len(jwtKey) == 0 {
		return 0, "", "", errors.New("jwt Key not initialized")
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return 0, "", "", err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", "", errors.New("invalid claims type")
	}

	if claims.UserID == 0 || claims.Email == "" {
		return 0, "", "", errors.New("invalid user claims")
	}

	return claims.UserID, claims.Email, claims.Role, nil
}
