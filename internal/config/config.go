package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (*Config, error) {
	// Проверяем, есть ли файл
	if _, err := os.Stat(".env"); err != nil {
		return nil, fmt.Errorf(".env file not found: %w", err)
	}

	// Загружаем .env
	err := godotenv.Load("C:/Users/ilyas/GolandProjects/onlineSubscriptions/.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	// Проверяем, что переменные не пустые
	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			return nil, fmt.Errorf("missing required env variable: %s", varName)
		}
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	fmt.Printf("Config loaded: %+v\n", cfg) // Вывод всех полей
	return cfg, nil
}
