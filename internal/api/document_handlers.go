package api

import (
	"encoding/json"
	"net/http"

	"github.com/IlyaChern12/rtce/internal/models"
	"github.com/IlyaChern12/rtce/internal/repository"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	repo *repository.DocumentRepository
}

type Request struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewDocumentHanler(repo *repository.DocumentRepository) *DocumentHandler {
	return &DocumentHandler{
		repo: repo,
	}
}

// /POST - создание дока
func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	doc := &models.Document{
		ID:     uuid.New().String(),
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
	}

	if err := h.repo.Create(r.Context(), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": doc.ID}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}