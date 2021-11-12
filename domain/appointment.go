package domain

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentStatus string

var (
	Confirmed            AppointmentStatus = "confirmed"
	AwaitingConfirmation AppointmentStatus = "awaiting confirmation"
)

type Appointment struct {
	ID                uuid.UUID       `json:"id" gorm:"primary_key"`
	TreatmentCenterID uuid.UUID       `json:"treatment_center_id"`
	TreatmentCenter   TreatmentCenter `json:"-" binding:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	StartTime         time.Time       `json:"start_time" binding:"required"`
	CreatedAt         *time.Time      `json:"created_at"`
	UpdatedAt         *time.Time      `json:"updated_at"`
}

type AppointmentBooking struct {
	ID            uuid.UUID         `json:"id"`
	AppointmentID uuid.UUID         `json:"appointment_id" gorm:"association_foreignkey:ID"`
	PatientID     uuid.UUID         `json:"patient_id" gorm:"association_foreignkey:ID"`
	Status        AppointmentStatus `json:"status"`
}
