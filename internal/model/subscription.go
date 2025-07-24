package model

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ServiceName string     `json:"service_name" gorm:"not null"`
	Price       int        `json:"price" gorm:"not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
