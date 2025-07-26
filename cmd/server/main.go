// @title Online Subscriptions API
// @version 1.0
// @description API для управления онлайн-подписками пользователей
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	_ "github.com/IlyaStarshinov/onlineSubscriptions/docs"
	"github.com/IlyaStarshinov/onlineSubscriptions/internal/repository"

	"github.com/IlyaStarshinov/onlineSubscriptions/internal/handler"
)

func main() {
	db, err := repository.InitRepository()
	if err != nil {
		log.Fatalf("failed to initialize database, got error %v", err)
	}

	router := handler.SetupRouter(db)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server, got error %v", err)
	}
}
