package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	errpkg "github.com/lekss361/curserv2/gateway/internal/errors"
	authsvc "github.com/lekss361/curserv2/gateway/internal/service"
)

type keyUser struct{}

func userKey() keyUser { return keyUser{} }

type AuthHandler struct {
	authService authsvc.AuthService
}

func NewAuthHandler(authService authsvc.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.authService.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		login, ok, err := h.authService.ValidateToken(r.Context(), parts[1])
		if err != nil {
			http.Error(w, "token validation error", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey(), login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AuthHandler) ErrorResponder(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errpkg.NotFoundError:
		http.Error(w, e.Error(), e.StatusCode())
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
