package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID                  uuid.UUID             `json:"id" gorm:"primary_key"`
	Email               string                `json:"email" binding:"required"`
	FirstName           string                `json:"first_name" binding:"required"`
	LastName            string                `json:"last_name" binding:"required"`
	AppointmentBookings []*AppointmentBooking `json:"_" binding:"-"`
	CreatedAt           *time.Time            `json:"created_at"`
	UpdatedAt           *time.Time            `json:"updated_at"`
}
