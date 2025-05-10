package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID            string   `json:"sub"`
	TokenType         string   `json:"typ"`
	UserType          UserType `json:"user"`
	PasswordChangedAt int64    `json:"lpc"`
	jwt.RegisteredClaims
}

type TokenService struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

func NewTokenService(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *TokenService {
	return &TokenService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *TokenService) GenerateAccessToken(userID string, userType UserType, passwordChangedAt *time.Time) (string, error) {
	if !userType.IsValid() {
		return "", fmt.Errorf("invalid user type: %s", userType)
	}
	return s.generateToken(userID, s.accessSecret, s.accessExpiry, TokenTypeAccess, userType, passwordChangedAt)
}

func (s *TokenService) GenerateRefreshToken(userID string, userType UserType, passwordChangedAt *time.Time) (string, error) {
	if !userType.IsValid() {
		return "", fmt.Errorf("invalid user type: %s", userType)
	}
	return s.generateToken(userID, s.refreshSecret, s.refreshExpiry, TokenTypeRefresh, userType, passwordChangedAt)
}

func (s *TokenService) generateToken(userID, secret string, expiry time.Duration, tokenType string, userType UserType, passwordChangedAt *time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(expiry).Unix(),
		"iat":  time.Now().Unix(),
		"typ":  tokenType,
		"user": string(userType),
		"lpc":  passwordChangedAt.UTC().UnixNano(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return s.validateToken(tokenString, s.accessSecret, TokenTypeAccess)
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return s.validateToken(tokenString, s.refreshSecret, TokenTypeRefresh)
}

func (s *TokenService) validateToken(tokenString, secret, expectedType string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, fmt.Errorf("invalid token type: expected %s", expectedType)
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
