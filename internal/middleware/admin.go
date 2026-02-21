package middleware

import "net/http"

// RequireAdmin rejects non-admin requests.
// Must be used after Authenticate middleware.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetRole(r.Context())
		if !ok || role != "admin" {
			respondForbidden(w, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
