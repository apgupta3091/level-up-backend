package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/anujgupta/level-up-backend/internal/auth"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyRole   contextKey = "role"
)

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return id, ok
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyRole).(string)
	return role, ok
}

// Authenticate validates the Bearer JWT and injects userID + role into context.
func Authenticate(authSvc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				respondUnauthorized(w, "missing or invalid authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ParseToken(tokenStr)
			if err != nil {
				respondUnauthorized(w, "invalid or expired token")
				return
			}

			if claims.TokenType != auth.TokenTypeAccess {
				respondUnauthorized(w, "access token required")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
