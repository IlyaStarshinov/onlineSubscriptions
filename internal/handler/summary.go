package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
)

// @Summary Получить сумму подписок
// @Description Выводит общую сумму подписок по фильтрам
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Начало периода (MM-YYYY)"
// @Param end_date query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error"
// @Router /subscriptions/summary [get]
func (h *Handler) GetSubscriptionSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		respondError(w, http.StatusBadRequest, "Missing start date or end date")
		return
	}
	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start date")
		return
	}
	endDate, err := time.Parse("01-2006", endDateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end date")
		return
	}
	var userUUID uuid.UUID
	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}
	}
	var total int
	query := h.DB.Model(&model.Subscription{}).
		Where("start_date >= ? AND start_date <= ?", startDate, endDate)
	if userID != "" {
		query = query.Where("user_id = ?", userUUID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if err := query.Select("SUM(price)").Scan(&total).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch subscription summary")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"total_price": total})
}
