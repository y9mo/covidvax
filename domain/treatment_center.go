package domain

import (
	"time"

	"github.com/google/uuid"
)

type TreatmentCenter struct {
	ID           uuid.UUID      `json:"id" gorm:"primary_key"`
	Name         string         `json:"name" binding:"required"`
	Address      string         `json:"address" binding:"required"`
	Phone        string         `json:"phone" binding:"required"`
	Appointments []*Appointment `json:"-" binding:"-"`
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
}
