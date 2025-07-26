package handler

import (
	"net/http"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
