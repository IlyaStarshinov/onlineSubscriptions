package main

import (
	"github.com/IlyaStarshinov/onlineSubscriptions/internal/repository"
	"log"
	"net/http"

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
