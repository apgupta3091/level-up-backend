package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

const stripeRawBodyKey contextKey = "stripeRawBody"

func GetStripeRawBody(ctx context.Context) ([]byte, bool) {
	b, ok := ctx.Value(stripeRawBodyKey).([]byte)
	return b, ok
}

// StripeRawBody reads the request body before any parsing and stores the raw
// bytes in the context. This is required for Stripe webhook signature verification.
func StripeRawBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const maxBodySize = 65 << 10 // 65 KB
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		ctx := context.WithValue(r.Context(), stripeRawBodyKey, rawBody)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
