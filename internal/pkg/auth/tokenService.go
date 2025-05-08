package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserType string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
	UserTypeRider UserType = "rider"
)

func (u UserType) IsValid() bool {
	switch u {
	case UserTypeAdmin, UserTypeUser, UserTypeRider:
		return true
	}
	return false
}

type TokenClaims struct {
	UserID    string   `json:"sub"`
	TokenType string   `json:"typ"`
	UserType  UserType `json:"user"`
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

func (s *TokenService) GenerateAccessToken(userID string, userType UserType) (string, error) {
	if !userType.IsValid() {
		return "", fmt.Errorf("invalid user type: %s", userType)
	}
	return s.generateToken(userID, s.accessSecret, s.accessExpiry, TokenTypeAccess, userType)
}

func (s *TokenService) GenerateRefreshToken(userID string, userType UserType) (string, error) {
	if !userType.IsValid() {
		return "", fmt.Errorf("invalid user type: %s", userType)
	}
	return s.generateToken(userID, s.refreshSecret, s.refreshExpiry, TokenTypeRefresh, userType)
}

func (s *TokenService) generateToken(userID, secret string, expiry time.Duration, tokenType string, userType UserType) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(expiry).Unix(),
		"iat":  time.Now().Unix(),
		"typ":  tokenType,
		"user": string(userType),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
