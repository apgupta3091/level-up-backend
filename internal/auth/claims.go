package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID `json:"uid"`
	Role      string    `json:"role"`
	TokenType string    `json:"typ"`
}
