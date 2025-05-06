package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authclient "github.com/lekss361/curserv2/gateway/internal/clients/auth"
	"github.com/lekss361/curserv2/gateway/internal/config"
	handlerpkg "github.com/lekss361/curserv2/gateway/internal/handler"
	repopkg "github.com/lekss361/curserv2/gateway/internal/repository"
	servicepkg "github.com/lekss361/curserv2/gateway/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	authCli := authclient.NewClient(cfg.AuthServiceURL)
	userRepo := repopkg.NewInMemoryUserRepo()

	authSvc := servicepkg.NewAuthService(authCli, userRepo)
	currSvc, err := servicepkg.NewCurrencyService(cfg.CurrencyServiceURL)
	if err != nil {
		log.Fatalf("currency service init: %v", err)
	}

	authHandler := handlerpkg.NewAuthHandler(authSvc)
	currencyHandler := handlerpkg.NewCurrencyHandler(currSvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	r.Post("/login", authHandler.Login)

	r.Route("/currency", func(r chi.Router) {
		r.Use(authHandler.ValidateToken)

		r.Get("/{date}", currencyHandler.Get)
		r.Get("/history", currencyHandler.History)
	})

	log.Printf("Gateway listening on %s", cfg.Server.BindAddr)
	if err := http.ListenAndServe(cfg.Server.BindAddr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
