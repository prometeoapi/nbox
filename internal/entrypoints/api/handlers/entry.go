package handlers

import (
	"encoding/json"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"nbox/internal/entrypoints/api/response"
	"nbox/internal/usecases"
	"net/http"
)

type EntryHandler struct {
	entryAdapter domain.EntryAdapter
	entryUseCase *usecases.EntryUseCase
}

func NewEntryHandler(entryAdapter domain.EntryAdapter, entryUseCase *usecases.EntryUseCase) *EntryHandler {
	return &EntryHandler{entryAdapter: entryAdapter, entryUseCase: entryUseCase}
}

func (h *EntryHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var entries []models.Entry

	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	//results := h.entryUseCase.Upsert(ctx, entries)
	//err := h.entryAdapter.Upsert(ctx, entries)
	//if err != nil {
	//	response.Error(w, r, err, http.StatusBadRequest)
	//	return
	//}
	response.Success(w, r, h.entryUseCase.Upsert(ctx, entries))
}

func (h *EntryHandler) ListByPrefix(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prefix := r.URL.Query().Get("v")

	entries, err := h.entryAdapter.List(ctx, prefix)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	response.Success(w, r, entries)
}

func (h *EntryHandler) GetByKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.URL.Query().Get("v")
	entry, err := h.entryAdapter.Retrieve(ctx, key)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	response.Success(w, r, entry)
}

func (h *EntryHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.URL.Query().Get("v")
	err := h.entryAdapter.Delete(ctx, key)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	response.Success(w, r, map[string]string{"message": "ok"})
}

func (h *EntryHandler) Tracking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := r.URL.Query().Get("v")

	entries, err := h.entryAdapter.Tracking(ctx, key)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	response.Success(w, r, entries)
}
