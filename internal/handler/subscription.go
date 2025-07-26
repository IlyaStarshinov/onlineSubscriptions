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

type CreateSubscriptionInput struct {
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

// @Summary Создать подписку
// @Description Создаёт новую онлайн-подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body handler.CreateSubscriptionInput true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {string} string "Bad Request"
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input CreateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if input.ServiceName == "" {
		log.Printf("Validation error: service name is empty")
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	if input.Price < 0 {
		log.Printf("Price is negative")
		http.Error(w, "Price must be > 0", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Printf("Failed to parse user UUID: %v", err)
		http.Error(w, "user_id must be valid UUID", http.StatusBadRequest)
		return
	}

	starDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		log.Printf("Failed to parse start date: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var endDatePrt *time.Time
	if input.EndDate != nil {
		endDate, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			log.Printf("Failed to parse end date: %v", err)
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
		log.Printf("Failed to create subscription: %v", err)
		http.Error(w, fmt.Sprintf("failed to create subscription: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("created subscription: service=%s, user_id=%s, price=%d", sub.ServiceName, userUUID.String(), sub.Price)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// @Summary Получить все подписки
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 400 {string} string "Bad Request"
// @Router /subscriptions [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	var sub []model.Subscription

	if err := h.DB.Find(&sub).Error; err != nil {
		log.Printf("Failed to find subscriptions: %v", err)
		http.Error(w, "Failed to fetch subscriptions", http.StatusBadRequest)
		return
	}

	log.Printf("found subscription: %v", sub)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// @Summary Получить подписки по user_id
// @Description Возвращает список подписок определенного пользователя
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Success 200 {array} model.Subscription
// @Failure 400 {string} string "Bad Request"
// @Router /subscriptions/{user_id} [get]
func (h *Handler) GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("Failed to parse user UUID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var subs []model.Subscription
	if err := h.DB.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		log.Printf("Failed to find subscriptions: %v", err)
		http.Error(w, "Failed to fetch subscriptions", http.StatusBadRequest)
		return
	}

	log.Printf("found subscriptions: %v", subs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Param id path string true "UUID подписки"
// @Success 204
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subIDStr := vars["id"]

	subId, err := uuid.Parse(subIDStr)
	if err != nil {
		log.Printf("Failed to parse user UUID: %v", err)
		http.Error(w, "Invalid sub ID", http.StatusBadRequest)
		return
	}

	result := h.DB.Where("id = ?", subId).Delete(&model.Subscription{})
	if result.Error != nil {
		log.Printf("Failed to delete subscription: %v", result.Error)
		http.Error(w, "Failed to delete subscription", http.StatusBadRequest)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("No subscription found to delete with id: %s", subId)
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	log.Printf("deleted subscription: %v", subId)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить сумму подписок
// @Description Выводит общую сумму подписок по фильтрам (по пользователю, по сервису, по периоду)
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Начало периода (MM-YYYY)"
// @Param end_date query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /subscriptions/summary [get]
func (h *Handler) GetSubscriptionSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		log.Printf("Start date and end date are required")
		http.Error(w, "Missing start date or end date", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		log.Printf("Failed to parse start date: %v", err)
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("01-2006", endDateStr)
	if err != nil {
		log.Printf("Failed to parse end date: %v", err)
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}
	var userUUID uuid.UUID
	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err != nil {
			log.Printf("Failed to parse user UUID: %v", err)
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
		log.Printf("Failed to fetch subscription summary: %v", err)
		http.Error(w, "Failed to fetch subscription summary", http.StatusInternalServerError)
		return
	}

	log.Printf("Subscription summary: user_id=%s, service_name=%s, sum=%d", userUUID, serviceName, sum)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{
		"total_price": sum,
	})
}
