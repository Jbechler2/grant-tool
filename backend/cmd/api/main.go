package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jbechler2/grant-tool/backend/config"
	"github.com/jbechler2/grant-tool/backend/internal/auth"
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
	refreshTokenService := service.NewRefreshTokenService(database, queries)
	authService := service.NewAuthService(queries, cfg.JWTSecret, cfg.JWTExpiryMinutes, refreshTokenService)
	authHandler := handler.NewAuthHandler(authService, cfg.IsProduction)
	clientService := service.NewClientService(queries)
	clientHandler := handler.NewClientHandler(clientService)
	grantService := service.NewGrantService(queries)
	grantHandler := handler.NewGrantHandler(grantService)
	applicationService := service.NewApplicationService(queries)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.NewJWTMiddleware(cfg.JWTSecret))
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/clients", clientHandler.CreateClient)
			r.Get("/clients", clientHandler.GetAllClients)
			r.Get("/clients/{id}", clientHandler.GetClientByID)
			r.Put("/clients/{id}", clientHandler.UpdateClient)
			r.Delete("/clients/{id}", clientHandler.DeleteClient)
			r.Get("/clients/{id}/applications", applicationHandler.GetAllApplicationsByClientID)

			r.Post("/grants", grantHandler.CreateGrant)
			r.Get("/grants", grantHandler.GetAllGrants)
			r.Get("/grants/{id}", grantHandler.GetGrantByID)
			r.Put("/grants/{id}", grantHandler.UpdateGrant)
			r.Delete("/grants/{id}", grantHandler.DeleteGrant)
			r.Get("/grants/{id}/deadlines", grantHandler.GetDeadlinesByGrantID)
			r.Post("/grants/{id}/deadlines", grantHandler.AddDeadline)
			r.Delete("/grants/{id}/deadlines/{deadlineID}", grantHandler.DeleteDeadline)

			r.Post("/applications", applicationHandler.CreateApplication)
			r.Get("/applications", applicationHandler.GetAllApplicationsByUserID)
			r.Get("/applications/{id}", applicationHandler.GetApplicationByID)
			r.Put("/applications/{id}", applicationHandler.UpdateApplication)
			r.Put("/applications/{id}/publish", applicationHandler.PublishApplication)
			r.Delete("/applications/{id}", applicationHandler.DeleteApplication)
		})
	})

	log.Println("grant-tool API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
