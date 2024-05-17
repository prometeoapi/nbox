package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"time"
)

type Api struct {
	Engine http.Handler
}

func NewApi(box *BoxHandler, entry *EntryHandler) *Api {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(response.NotFound)
	r.MethodNotAllowed(response.MethodNotAllowed)

	r.Post("/api/box", box.UpsertBox)
	r.Head("/api/box/{service}/{stage}/{template}", box.Exist)
	r.Get("/api/box/{service}/{stage}/{template}", box.Retrieve)
	r.Get("/api/box/{service}/{stage}/{template}/build", box.Build)
	r.Post("/api/entry", entry.Upsert)
	r.Get("/api/entry/key", entry.GetByKey)
	r.Get("/api/entry/prefix", entry.ListByPrefix)

	return &Api{
		Engine: r,
	}
}
