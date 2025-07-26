package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UpdateSubscriptionInput входные данные для обновления подписки
type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       *int    `json:"price,omitempty" example:"399"`
	StartDate   *string `json:"start_date,omitempty" example:"02-2023"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2023"`
}

// @Summary Обновить подписку
// @Description Обновляет поля существующей подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Param input body handler.UpdateSubscriptionInput true "Данные для обновления"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
// @Failure 404 {object} handler.ErrorResponse "Not Found"
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	subIDStr := mux.Vars(r)["id"]
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var input UpdateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var sub model.Subscription
	if err := h.DB.First(&sub, "id = ?", subID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	if input.ServiceName != nil {
		sub.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		if *input.Price < 0 {
			respondError(w, http.StatusBadRequest, "Price must be > 0")
			return
		}
		sub.Price = *input.Price
	}
	if input.StartDate != nil {
		t, err := time.Parse("01-2006", *input.StartDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid start date format")
			return
		}
		sub.StartDate = t
	}
	if input.EndDate != nil {
		if *input.EndDate == "" {
			sub.EndDate = nil
		} else {
			t, err := time.Parse("01-2006", *input.EndDate)
			if err != nil {
				respondError(w, http.StatusBadRequest, "Invalid end date format")
				return
			}
			sub.EndDate = &t
		}
	}

	if err := h.DB.Save(&sub).Error; err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to update subscription: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}
