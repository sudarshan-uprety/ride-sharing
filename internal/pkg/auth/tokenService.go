package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func (s *TokenService) GenerateAccessToken(userID string, passwordChangedDT *time.Time) (string, error) {
	return s.generateToken(userID, s.accessSecret, s.accessExpiry, passwordChangedDT, TokenTypeAccess)
}

func (s *TokenService) GenerateRefreshToken(userID string, passwordChangedDT *time.Time) (string, error) {
	return s.generateToken(userID, s.refreshSecret, s.refreshExpiry, passwordChangedDT, TokenTypeRefresh)
}

func (s *TokenService) generateToken(userID, secret string, expiry time.Duration, passwordChangedDT *time.Time, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
		"iat": time.Now().Unix(),
		"typ": tokenType,
		"pca": passwordChangedDT.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
