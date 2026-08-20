package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"klikku/internal/config"
)

type Claims struct {
	UserID     string `json:"user_id"`
	MerchantID string `json:"merchant_id"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, merchantID, role, email string, cfg *config.Config) (string, error) {
	claims := Claims{
		UserID:     userID,
		MerchantID: merchantID,
		Role:       role,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "klikku",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func GenerateRefreshToken(userID string, cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "klikku",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
