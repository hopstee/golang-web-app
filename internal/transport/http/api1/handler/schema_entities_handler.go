package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/service"
	"mobile-backend-boilerplate/pkg/helper/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SchemaEntityHandler struct {
	schemaEntitiesService *service.SchemaEntityService
}

func NewSchemaEntityHandler(schemaEntitiesService *service.SchemaEntityService) *SchemaEntityHandler {
	return &SchemaEntityHandler{
		schemaEntitiesService: schemaEntitiesService,
	}
}

func (h *SchemaEntityHandler) GetEntityNamesByType(w http.ResponseWriter, r *http.Request) {
	entityType := chi.URLParam(r, "type")

	names, err := h.schemaEntitiesService.GetEntitiesName(r.Context(), entityType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, names)
}

func (h *SchemaEntityHandler) GetEntitySchema(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	entityType := chi.URLParam(r, "type")
	schema, err := h.schemaEntitiesService.GetEntitySchema(r.Context(), entityType, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, schema)
}

func (h *SchemaEntityHandler) GetEntityData(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	data, err := h.schemaEntitiesService.GetEntityData(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, data)
}

func (h *SchemaEntityHandler) UpdateEntityData(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.schemaEntitiesService.UpdateEntityData(r.Context(), slug, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
