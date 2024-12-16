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

// Upsert
// @Summary Upsert entries
// @Description insert / update vars
// @Tags entry
// @Accept json
// @Produce json
// @Param data body []models.Entry true "Upsert template"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} map[string]string ""
// @Failure 400 {object} problem.ProblemDetail "Bad Request"
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/entry [post]
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

// ListByPrefix
// @Summary Filter by prefix
// @Description list all keys by path
// @Tags entry
// @Produce json
// @Param v query string true "key path"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} []models.Entry ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/entry/prefix [get]
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

// GetByKey
// @Summary Retrieve key
// @Description detail
// @Tags entry
// @Produce json
// @Param v query string true "key path"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} models.Entry ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/entry/key [get]
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

// DeleteKey
// @Summary Delete
// @Description delete keys & children
// @Tags entry
// @Produce json
// @Param v query string true "key path"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} object{message=string} ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/entry/key [delete]
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

// Tracking
// @Summary History
// @Description history changes
// @Tags entry
// @Produce json
// @Param v query string true "key path"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} []models.Tracking ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/track/key [get]
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
