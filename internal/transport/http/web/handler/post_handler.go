package handler

import (
	"encoding/json"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// TODO: переписать хэндлер для работы с templ + htmx

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) GetAllPublic(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPublicPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPublicPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	post := repository.Post{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: r.FormValue("is_public") == "true",
	}

	_, fileHeader, _ := r.FormFile("hero_img")

	id, err := h.postService.CreatePost(post, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(id)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)

	post := repository.Post{
		ID:       id,
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: r.FormValue("is_public") == "true",
	}

	_, fileHeader, _ := r.FormFile("hero_img")

	err = h.postService.UpdatePost(post, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
