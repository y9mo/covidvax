package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/y9mo/covidvax/domain"
	"go.uber.org/zap"
)

type Appointments interface {
	Create(appointment *domain.Appointment) error
	Update(appointment *domain.Appointment) error
	Delete(appointment *domain.Appointment) error
	FindByID(id uuid.UUID) (*domain.Appointment, error)
	All() (result []*domain.Appointment, err error)
	AllByTreatmentCenterID(treatmentCenterID uuid.UUID) (result []*domain.Appointment, err error)
	AllAvailable() (result []*domain.Appointment, err error)
	AllBookedByTreatmentCenterIDForDate(treatmentCenterID uuid.UUID, date time.Time) (result []*domain.Appointment, err error)
}

type appointments struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAppointments(db *gorm.DB, logger *zap.Logger) Appointments {
	return appointments{db: db, logger: logger}
}

func (r appointments) Create(appointment *domain.Appointment) error {
	err := r.db.Debug().Create(appointment).Error
	return handleGormError(err, r.logger)
}

func (r appointments) Update(appointment *domain.Appointment) error {
	err := r.db.Debug().Update(appointment).Error
	return handleGormError(err, r.logger)
}

func (r appointments) Delete(appointment *domain.Appointment) error {
	err := r.db.Debug().Delete(appointment).Error
	return handleGormError(err, r.logger)
}

func (r appointments) FindByID(id uuid.UUID) (*domain.Appointment, error) {
	appointment := domain.Appointment{}
	err := r.db.Debug().Where("id = ?", id).Find(&appointment).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r appointments) All() (result []*domain.Appointment, err error) {
	err = r.db.Debug().Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r appointments) AllByTreatmentCenterID(treatmentCenterID uuid.UUID) (result []*domain.Appointment, err error) {
	err = r.db.Debug().Where("treatment_center_id = ?", treatmentCenterID).Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r appointments) AllAvailable() (result []*domain.Appointment, err error) {
	err = r.db.Debug().Joins("LEFT JOIN appointment_bookings ab on appointments.id = ab.appointment_id").
		Where("ab.appointment_id is NULL").Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r appointments) AllBookedByTreatmentCenterIDForDate(treatmentCenterID uuid.UUID,
	date time.Time) (result []*domain.Appointment, err error) {
	err = r.db.Debug().Joins("LEFT JOIN appointment_bookings ab on appointments.id = ab.appointment_id").
		Where("ab.appointment_id is NOT NULL AND treatment_center_id = ?", treatmentCenterID).
		Where("date_trunc('day', appointments.start_time) = date_trunc('day', ?::timestamptz)", date).
		Find(&result).Error

	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}
