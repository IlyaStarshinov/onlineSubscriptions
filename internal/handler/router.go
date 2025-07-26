package handler

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	_ "github.com/IlyaStarshinov/onlineSubscriptions/docs"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	h := NewHandler(db)
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/subscriptions", h.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscriptions", h.GetSubscription).Methods("GET")
	r.HandleFunc("/subscriptions/{user_id}", h.GetSubscriptionsByUserID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", h.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/subscriptions/summary", h.GetSubscriptionSummary).Methods("GET")

	return r
}
