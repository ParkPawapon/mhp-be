package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	Role      constants.Role `json:"role"`
	TokenType TokenType      `json:"type"`
	SessionID string         `json:"sid"`
	jwt.RegisteredClaims
}

func NewAccessToken(userID uuid.UUID, role constants.Role, sessionID uuid.UUID, cfg config.JWTConfig) (string, error) {
	return buildToken(userID, role, sessionID, TokenTypeAccess, cfg.AccessTTL, cfg)
}

func NewRefreshToken(userID uuid.UUID, role constants.Role, sessionID uuid.UUID, cfg config.JWTConfig) (string, error) {
	return buildToken(userID, role, sessionID, TokenTypeRefresh, cfg.RefreshTTL, cfg)
}

func ParseToken(tokenString string, cfg config.JWTConfig) (*Claims, error) {
	claims := &Claims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	_, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func buildToken(userID uuid.UUID, role constants.Role, sessionID uuid.UUID, tokenType TokenType, ttl time.Duration, cfg config.JWTConfig) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Role:      role,
		TokenType: tokenType,
		SessionID: sessionID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
