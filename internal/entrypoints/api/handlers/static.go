package handlers

import (
	"nbox/internal/application"
	"nbox/internal/entrypoints/api/response"
	"net/http"
)

type StaticHandler struct {
	config *application.Config
}

func NewStaticHandler(config *application.Config) *StaticHandler {
	return &StaticHandler{
		config: config,
	}
}

func (s *StaticHandler) Environments(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, s.config.AllowedPrefixes)
}
