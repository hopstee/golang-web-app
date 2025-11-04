package handler

import (
	"mobile-backend-boilerplate/internal/filestorage"
	"net/http"
	"time"
)

type FilesHandler struct {
	fileStorage filestorage.FileStorage
}

func NewFilesHandler(fileStorage filestorage.FileStorage) *FilesHandler {
	return &FilesHandler{
		fileStorage: fileStorage,
	}
}

func (h *FilesHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData := make([]byte, handler.Size)
	_, err = file.Read(fileData)
	if err != nil {
		http.Error(w, "Could not read uploaded file", http.StatusInternalServerError)
		return
	}

	fileName := time.Now().Format("20060102150405") + "_" + handler.Filename
	path, err := h.fileStorage.SaveFile(fileName, fileData)
	if err != nil {

		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(path))
}

func (h *FilesHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	err := h.fileStorage.DeleteFile(filePath)
	if err != nil {
		http.Error(w, "Could not delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
