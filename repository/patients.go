package repository

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/y9mo/covidvax/domain"
	"go.uber.org/zap"
)

type Patients interface {
	Create(patient *domain.Patient) error
	Update(patient *domain.Patient) error
	FindByID(id uuid.UUID) (*domain.Patient, error)
	Delete(patient *domain.Patient) error
	All() ([]*domain.Patient, error)
}
type patients struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPatients(db *gorm.DB, logger *zap.Logger) Patients {
	return patients{
		db:     db,
		logger: logger,
	}
}

func (r patients) Create(patient *domain.Patient) error {
	err := r.db.Debug().Create(patient).Error
	return handleGormError(err, r.logger)
}

func (r patients) Update(patient *domain.Patient) error {
	err := r.db.Debug().Update(patient).Error
	return handleGormError(err, r.logger)
}

func (r patients) FindByID(id uuid.UUID) (*domain.Patient, error) {
	patient := domain.Patient{}
	err := r.db.Debug().Where("id = ?", id).Find(&patient).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r patients) Delete(patient *domain.Patient) error {
	err := r.db.Debug().Delete(patient).Error
	return handleGormError(err, r.logger)
}

func (r patients) All() (result []*domain.Patient, err error) {
	err = r.db.Debug().Find(&result).Error
	err = handleGormError(err, r.logger)
	if err != nil {
		return nil, err
	}
	return result, nil
}
