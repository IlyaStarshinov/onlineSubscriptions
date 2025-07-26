package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
)

// CreateSubscriptionInput входные данные для создания подписки
type CreateSubscriptionInput struct {
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"599"`
	UserID      string  `json:"user_id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`
	StartDate   string  `json:"start_date" example:"01-2023"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2023"`
}

// @Summary Создать подписку
// @Description Создаёт новую онлайн-подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body handler.CreateSubscriptionInput true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} handler.ErrorResponse
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input CreateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if input.ServiceName == "" {
		respondError(w, http.StatusBadRequest, "Service name is required")
		return
	}
	if input.Price < 0 {
		respondError(w, http.StatusBadRequest, "Price must be > 0")
		return
	}
	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "user_id must be valid UUID")
		return
	}
	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start date format")
		return
	}
	var endDatePtr *time.Time
	if input.EndDate != nil {
		t, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid end date format")
			return
		}
		endDatePtr = &t
	}
	sub := model.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}
	if err := h.DB.Create(&sub).Error; err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("failed to create subscription: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}
