package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/pkg/helper/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PagesHandler struct {
	pagesService *service.PagesService
}

func NewPagesHandler(pagesService *service.PagesService) *PagesHandler {
	return &PagesHandler{
		pagesService: pagesService,
	}
}

func (h *PagesHandler) GetAllPagesSchemas(w http.ResponseWriter, r *http.Request) {
	schemas, err := h.pagesService.GetAllPagesSchemas(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, schemas)
}

func (h *PagesHandler) GetPageSchema(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	schema, err := h.pagesService.GetPageSchema(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, schema)
}

func (h *PagesHandler) GetPageData(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	schema, data, err := h.pagesService.GetPageData(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, map[string]interface{}{
		"schema": schema,
		"data":   data,
	})
}

func (h *PagesHandler) UpdatePageData(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.pagesService.UpdatePageData(r.Context(), slug, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PagesHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if err := h.pagesService.DeletePageData(r.Context(), slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
