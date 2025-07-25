package handler

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	h := NewHandler(db)
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/subscriptions", h.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscriptions", h.GetSubscription).Methods("GET")
	r.HandleFunc("/subscriptions/{user_id}", h.GetSubscriptionsByUserID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", h.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/subscriptions/summary", h.GetSubscriptionSummary).Methods("GET")

	return r
}
