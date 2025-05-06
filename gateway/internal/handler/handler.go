package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authsvc "github.com/lekss361/curserv2/gateway/internal/service"
	currsvc "github.com/lekss361/curserv2/gateway/internal/service"
)

type Handler struct {
	Auth     *AuthHandler
	Currency *CurrencyHandler
}

func NewHandler(authSvc authsvc.AuthService, currSvc currsvc.CurrencyService) *Handler {
	return &Handler{
		Auth:     NewAuthHandler(authSvc),
		Currency: NewCurrencyHandler(currSvc),
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	r.Post("/login", h.Auth.Login)

	r.Route("/currency", func(r chi.Router) {
		r.Use(h.Auth.ValidateToken)

		r.Get("/{date}", h.Currency.Get)

		r.Get("/history", h.Currency.History)
	})

	return r
}
