package handlers

import (
	"encoding/json"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"nbox/internal/entrypoints/api/response"
	"net/http"
)

type EntryHandler struct {
	entryAdapter domain.EntryAdapter
}

func NewEntryHandler(entryAdapter domain.EntryAdapter) *EntryHandler {
	return &EntryHandler{entryAdapter: entryAdapter}
}

func (h *EntryHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var entries []models.Entry

	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}

	err := h.entryAdapter.Upsert(ctx, entries)
	if err != nil {
		response.Error(w, r, err, http.StatusBadRequest)
		return
	}
	response.Success(w, r, map[string]string{"message": "ok"})
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
