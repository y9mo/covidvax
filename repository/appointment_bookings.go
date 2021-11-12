package repository

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/y9mo/covidvax/domain"
	"go.uber.org/zap"
)

type AppointmentBookings interface {
	Create(appointmentBooking *domain.AppointmentBooking) error
	Update(appointmentBooking *domain.AppointmentBooking) error
	Delete(appointmentBooking *domain.AppointmentBooking) error
	FindByID(id uuid.UUID) (*domain.AppointmentBooking, error)
}

type appointmentBookings struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAppointmentBookings(db *gorm.DB, logger *zap.Logger) AppointmentBookings {
	return appointmentBookings{db: db, logger: logger}
}

func (r appointmentBookings) Create(appointmentBooking *domain.AppointmentBooking) error {
	err := r.db.Debug().Create(appointmentBooking).Error
	return handleGormError(err, r.logger)
}

func (r appointmentBookings) Update(appointmentBooking *domain.AppointmentBooking) error {
	err := r.db.Debug().Update(appointmentBooking).Error
	return handleGormError(err, r.logger)
}

func (r appointmentBookings) Delete(appointmentBooking *domain.AppointmentBooking) error {
	err := r.db.Debug().Delete(appointmentBooking).Error
	return handleGormError(err, r.logger)
}

func (r appointmentBookings) FindByID(id uuid.UUID) (*domain.AppointmentBooking, error) {
	appointmentBooking := domain.AppointmentBooking{}
	err := r.db.Debug().Where("id = ?", id).Find(&appointmentBooking).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return &appointmentBooking, nil
}

func (r appointmentBookings) All() (result []*domain.AppointmentBooking, err error) {
	err = r.db.Debug().Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}
