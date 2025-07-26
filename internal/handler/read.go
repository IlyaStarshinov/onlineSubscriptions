package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary Получить все подписки
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
// @Router /subscriptions [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	var subs []model.Subscription
	if err := h.DB.Find(&subs).Error; err != nil {
		respondError(w, http.StatusBadRequest, "Failed to fetch subscriptions")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// @Summary Получить подписки по user_id
// @Description Возвращает список подписок определённого пользователя
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Success 200 {array} model.Subscription
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
// @Router /subscriptions/{user_id} [get]
func (h *Handler) GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var subs []model.Subscription
	if err := h.DB.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		respondError(w, http.StatusBadRequest, "Failed to fetch subscriptions")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}
