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
	"gorm.io/gorm"
)

// ErrorResponse — единый формат JSON‑ошибок
type ErrorResponse struct {
	// Описание ошибки
	Error string `json:"error" example:"описание ошибки"`
}

// CreateSubscriptionInput — входная модель для создания подписки
type CreateSubscriptionInput struct {
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"599"`
	UserID      string  `json:"user_id" example:"a1b2c3d4-e5f6-7g8h-9i0j-k1l2m3n4o5p6"`
	StartDate   string  `json:"start_date" example:"01-2023"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2023"`
}

// UpdateSubscriptionInput — входная модель для обновления подписки
type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

// Handler хранит зависимости (БД)
type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// @Summary Создать подписку
// @Description Создаёт новую онлайн-подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body handler.CreateSubscriptionInput true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
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

// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Param id path string true "UUID подписки"
// @Success 204
// @Failure 400 {object} handler.ErrorResponse "Bad Request"
// @Failure 404 {object} handler.ErrorResponse "Not Found"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	subIDStr := mux.Vars(r)["id"]
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}
	res := h.DB.Where("id = ?", subID).Delete(&model.Subscription{})
	if res.Error != nil {
		respondError(w, http.StatusBadRequest, "Failed to delete subscription")
		return
	}
	if res.RowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Subscription not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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

// вспомогательная функция для выдачи JSON‑ошибок
func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
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
