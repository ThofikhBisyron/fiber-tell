package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	accessSecret  string
	refreshSecret string
}

func NewJwtService(
	accessSecret string,
	refreshSecret string,
) *JwtService {
	return &JwtService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *JwtService) GenerateAccessToken(user_id int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": strconv.FormatInt(user_id, 10),
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(s.accessSecret))
}

func (s *JwtService) GenerateRefreshToken(user_id int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": strconv.FormatInt(user_id, 10),
		"type":    "refresh",
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(s.refreshSecret))
}

func (s *JwtService) ValidateAccessToken(
	tokenString string,
) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(s.accessSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims["type"] != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func (s *JwtService) ValidateRefreshToken(
	tokenString string,
) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(s.refreshSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims["type"] != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}
