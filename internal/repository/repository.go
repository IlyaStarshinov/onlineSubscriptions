package repository

import (
	"fmt"
	"github.com/IlyaStarshinov/onlineSubscriptions/internal/config"
	"github.com/IlyaStarshinov/onlineSubscriptions/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitRepository() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	fmt.Printf("Подключаюсь с параметрами:\nHost: %s\nPort: %s\nUser: %s\nDB: %s\n",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}

	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	log.Println("Database connected and migrated")
	return db, nil
}
