package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jbechler2/grant-tool/backend/config"
	"github.com/jbechler2/grant-tool/backend/internal/db"
	"github.com/jbechler2/grant-tool/backend/internal/handler"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBURL)
	defer database.Close()

	queries := repository.New(database)
	authService := service.NewAuthService(queries, cfg.JWTSecret, cfg.JWTExpiryMinutes)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	log.Println("grant-tool API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
