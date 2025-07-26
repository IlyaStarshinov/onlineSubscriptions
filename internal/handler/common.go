package handler

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// ErrorResponse формат ошибок API
type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}

// Handler базовый обработчик
type Handler struct {
	DB *gorm.DB
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// respondError отправляет ошибку в формате JSON
func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
