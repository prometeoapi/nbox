package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"nbox/internal/entrypoints/api/response"
	"nbox/internal/usecases"
	"net/http"
)

type BoxHandler struct {
	store      domain.StoreOperations
	boxUseCase *usecases.BoxUseCase
}

func NewBoxHandler(store domain.StoreOperations, boxUseCase *usecases.BoxUseCase) *BoxHandler {
	return &BoxHandler{store: store, boxUseCase: boxUseCase}
}

func (b *BoxHandler) UpsertBox(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	command := &models.Command[models.Box]{}
	if err := json.NewDecoder(r.Body).Decode(command); err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	result := b.store.UpsertBox(ctx, &command.Payload)
	response.Success(w, r, result)
}

func (b *BoxHandler) Exist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	service := chi.URLParam(r, "service")
	stage := chi.URLParam(r, "stage")
	template := chi.URLParam(r, "template")

	exists, err := b.store.BoxExists(ctx, service, stage, template)
	if err != nil {
		response.Error(w, r, err, http.StatusNotFound)
		return
	}

	response.Success(w, r, map[string]bool{"exist": exists})
}

func (b *BoxHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	service := chi.URLParam(r, "service")
	stage := chi.URLParam(r, "stage")
	template := chi.URLParam(r, "template")

	data, err := b.store.RetrieveBox(ctx, service, stage, template)
	if err != nil {
		response.Error(w, r, err, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(data)
}

func (b *BoxHandler) Build(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	service := chi.URLParam(r, "service")
	stage := chi.URLParam(r, "stage")
	template := chi.URLParam(r, "template")

	data, err := b.boxUseCase.BuildBox(ctx, service, stage, template)
	if err != nil {
		response.Error(w, r, err, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(data))
}
