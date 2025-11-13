package handler

import (
	"net/http"
	"path/filepath"
)

type AdminHandler struct {
	webappDir string
}

func NewAdminHandler(webappDir string) *AdminHandler {
	return &AdminHandler{
		webappDir: webappDir,
	}
}

func (h *AdminHandler) Handle(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(h.webappDir, "index.html")
	http.ServeFile(w, r, indexPath)
}
