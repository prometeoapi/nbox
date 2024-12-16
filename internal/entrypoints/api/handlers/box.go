package handlers

import (
	"encoding/json"
	"nbox/internal/domain"
	"nbox/internal/domain/models"
	"nbox/internal/entrypoints/api/response"
	"nbox/internal/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BoxHandler struct {
	store      domain.TemplateAdapter
	boxUseCase *usecases.BoxUseCase
}

type CommandBox struct {
	ID      string     `json:"id" example:"123"`
	Payload models.Box `json:"payload"`
}

func NewBoxHandler(store domain.TemplateAdapter, boxUseCase *usecases.BoxUseCase) *BoxHandler {
	return &BoxHandler{store: store, boxUseCase: boxUseCase}
}

// UpsertBox
// @Summary Upsert templates
// @Description insert or update templates on s3
// @Tags templates
// @Accept json
// @Produce json
// @Param data body CommandBox true "Upsert template"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} []string ""
// @Failure 400 {object} problem.ProblemDetail "Bad Request"
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box [post]
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

// Exist
// @Summary Exist template
// @Description Check the existence of the template
// @Tags templates
// @Accept json
// @Produce json
// @Param service path string true "service name"
// @Param stage path string true "stage"
// @Param template path string true "template name"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object}  object{exit=bool} ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box/{service}/{stage}/{template} [head]
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

// Retrieve
// @Summary Retrieve template
// @Description detail
// @Tags templates
// @Produce plain
// @Param service path string true "service name"
// @Param stage path string true "stage"
// @Param template path string true "template name"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object}  string ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box/{service}/{stage}/{template} [get]
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

// Build
// @Summary Build template
// @Description replace vars patterns
// @Tags templates
// @Produce plain
// @Param service path string true "service name"
// @Param stage path string true "stage"
// @Param template path string true "template name"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object}  string ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box/{service}/{stage}/{template}/build [get]
func (b *BoxHandler) Build(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	service := chi.URLParam(r, "service")
	stage := chi.URLParam(r, "stage")
	template := chi.URLParam(r, "template")
	args := make(map[string]string)

	for key := range r.URL.Query() {
		if key == "service" || key == "stage" || key == "template" {
			continue
		}
		args[key] = r.URL.Query().Get(key)
	}

	data, err := b.boxUseCase.BuildBox(ctx, service, stage, template, args)
	if err != nil {
		response.Error(w, r, err, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(data))
}

// List
// @Summary List templates
// @Description all templates
// @Tags templates
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object} []models.Box ""
// @Failure 400 {object} problem.ProblemDetail "Bad Request"
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box [get]
func (b *BoxHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := b.store.List(ctx)
	if err != nil {
		response.Error(w, r, err, http.StatusNotFound)
		return
	}
	response.Success(w, r, data)
}

// ListVars
// @Summary List vars template
// @Description show all vars in template
// @Tags templates
// @Produce json
// @Param service path string true "service name"
// @Param stage path string true "stage"
// @Param template path string true "template name"
// @Param authorization header string true "Bearer | Basic"
// @Success 200 {object}  []string ""
// @Failure 401 {object} problem.ProblemDetail "Unauthorized"
// @Failure 500 {object} problem.ProblemDetail "Internal error"
// @Router /api/box/{service}/{stage}/{template}/vars [get]
func (b *BoxHandler) ListVars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service := chi.URLParam(r, "service")
	stage := chi.URLParam(r, "stage")
	template := chi.URLParam(r, "template")

	data := b.boxUseCase.ListVars(ctx, service, stage, template)
	response.Success(w, r, data)
}
