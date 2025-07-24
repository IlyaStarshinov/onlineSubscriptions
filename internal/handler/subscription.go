package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type createSubscriptionInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input createSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if input.ServiceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	if input.Price < 0 {
		http.Error(w, "Price must be > 0", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		http.Error(w, "user_id must be valid UUID", http.StatusBadRequest)
		return
	}

	starDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var endDatePrt *time.Time
	if input.EndDate != nil {
		endDate, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		endDatePrt = &endDate
	}
	sub := model.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      userUUID,
		StartDate:   starDate,
		EndDate:     endDatePrt,
	}

	if err := h.DB.Create(&sub).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to create subscription: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	var sub []model.Subscription

	if err := h.DB.Find(&sub).Error; err != nil {
		http.Error(w, "Failed to fetch subscriptions", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var subs []model.Subscription
	if err := h.DB.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		http.Error(w, "Failed to fetch subscriptions", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subIDStr := vars["id"]

	subId, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, "Invalid sub ID", http.StatusBadRequest)
		return
	}

	result := h.DB.Where("id = ?", subId).Delete(&model.Subscription{})
	if result.Error != nil {
		http.Error(w, "Failed to delete subscription", http.StatusBadRequest)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSubscriptionSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing start date or end date", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("01-2006", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}
	var userUUID uuid.UUID
	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	}
	var sum int
	query := h.DB.Model(&model.Subscription{}).Where("start_date >= ? AND start_date <= ?", startDate, endDate)
	if userID != "" {
		query = query.Where("user_id = ?", userUUID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	err = query.Select("SUM(price)").Scan(&sum).Error
	if err != nil {
		http.Error(w, "Failed to fetch subscription summary", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{
		"total_price": sum,
	})

}
