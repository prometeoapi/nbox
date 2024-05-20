package api

import (
	"nbox/internal/entrypoints/api/health"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Api struct {
	Engine http.Handler
}

func NewApi(box *BoxHandler, entry *EntryHandler, healthCheck *health.Health) *Api {

	corsConfig := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Clientedni"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(healthCheck.Healthy("/health"))
	r.Use(corsConfig)

	r.NotFound(response.NotFound)
	r.MethodNotAllowed(response.MethodNotAllowed)

	r.Post("/api/box", box.UpsertBox)
	r.Get("/api/box", box.List)
	r.Head("/api/box/{service}/{stage}/{template}", box.Exist)
	r.Get("/api/box/{service}/{stage}/{template}", box.Retrieve)
	r.Get("/api/box/{service}/{stage}/{template}/build", box.Build)
	r.Post("/api/entry", entry.Upsert)
	r.Get("/api/entry/key", entry.GetByKey)
	r.Get("/api/entry/prefix", entry.ListByPrefix)
	r.Get("/health", entry.ListByPrefix)

	return &Api{
		Engine: r,
	}
}
