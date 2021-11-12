package repository

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/y9mo/covidvax/domain"
	"go.uber.org/zap"
)

type TreatmentCenters interface {
	Create(treatmentCenter *domain.TreatmentCenter) error
	Update(treatmentCenter *domain.TreatmentCenter) error
	FindByID(id uuid.UUID) (*domain.TreatmentCenter, error)
	Delete(treatmentCenter *domain.TreatmentCenter) error
	All() ([]*domain.TreatmentCenter, error)
}
type treatmentCenters struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTreatmentCenters(db *gorm.DB, logger *zap.Logger) TreatmentCenters {
	return treatmentCenters{
		db:     db,
		logger: logger,
	}
}

func (r treatmentCenters) Create(treatmentCenter *domain.TreatmentCenter) error {
	err := r.db.Debug().Create(treatmentCenter).Error
	return handleGormError(err, r.logger)
}

func (r treatmentCenters) Update(treatmentCenter *domain.TreatmentCenter) error {
	err := r.db.Debug().Update(treatmentCenter).Error
	return handleGormError(err, r.logger)
}

func (r treatmentCenters) FindByID(id uuid.UUID) (*domain.TreatmentCenter, error) {
	treatmentCenter := domain.TreatmentCenter{}
	err := r.db.Debug().Where("id = ?", id).Find(&treatmentCenter).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return &treatmentCenter, nil
}

func (r treatmentCenters) Delete(treatmentCenter *domain.TreatmentCenter) error {
	err := r.db.Debug().Delete(treatmentCenter).Error
	return handleGormError(err, r.logger)
}

func (r treatmentCenters) All() (result []*domain.TreatmentCenter, err error) {
	err = r.db.Debug().Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}
