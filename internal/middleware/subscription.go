package middleware

import (
	"net/http"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
)

// RequireActive rejects requests from users without an active subscription.
// Must be used after Authenticate middleware.
func RequireActive(queries *dbgen.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				respondUnauthorized(w, "not authenticated")
				return
			}

			status, err := queries.GetUserSubscriptionStatus(r.Context(), userID)
			if err != nil {
				respondForbidden(w, "could not verify subscription")
				return
			}

			if status != dbgen.SubscriptionStatusActive {
				respondForbidden(w, "active subscription required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
