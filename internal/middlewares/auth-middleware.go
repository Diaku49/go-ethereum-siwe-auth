package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const CtxWalletAddress contextKey = "wallet_address"

// For now this is a placeholder shape.
// Later you'll verify JWT and extract wallet address.
type TokenVerifier interface {
	Verify(token string) (walletAddress string, err error)
}

func RequireAuth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			addr, err := verifier.Verify(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxWalletAddress, addr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
