package middleware

import (
	"context"
	"net/http"
	"strings"

	authsvc "github.com/lekss361/curserv2/gateway/internal/service"
)

type ctxKeyUser struct{}

func UserFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(ctxKeyUser{}).(string)
	return u, ok
}

func AuthMiddleware(authService authsvc.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			login, ok, err := authService.ValidateToken(r.Context(), parts[1])
			if err != nil {
				http.Error(w, "token validation error", http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Add user login to context
			ctx := context.WithValue(r.Context(), ctxKeyUser{}, login)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
