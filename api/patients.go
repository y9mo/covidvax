package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/repository"
	"go.uber.org/zap"
)

type PatientsResponse struct {
	Patients []*domain.Patient `json:"patients,omitempty"`
	Error    *ErrorResponse    `json:"error,omitempty"`
}

type PatientResponse struct {
	Patient *domain.Patient `json:"patient,omitempty"`
	Error   *ErrorResponse  `json:"error,omitempty"`
}

type PatientsController struct {
	patientsRepository repository.Patients
	logger             *zap.Logger
}

type inputPatient struct {
	domain.Patient
	ID uuid.UUID `json:"-"`
}

func (c *inputPatient) buildModel() domain.Patient {
	return domain.Patient{
		ID: uuid.New(),
	}
}

func (c *inputPatient) updateModel(patient *domain.Patient) {
	// patient.Os = c.Os
}

func SetupPatient(
	router gin.IRouter,
	patientsRepository repository.Patients,
	logger *zap.Logger) {
	c := PatientsController{
		patientsRepository: patientsRepository,
		logger:             logger.With(zap.String("component", "PatientsController")),
	}
	g := router.Group("/patients")
	g.GET("/", c.IndexEndpoint)
	g.GET("/:patient_id", c.GetEndpoint)
	g.POST("/", c.CreateEndpoint)
	g.DELETE("/:patient_id", c.DeleteEndpoint)
	g.PUT("/:patient_id", c.UpdateEndpoint)
}

func extractPatientID(c *gin.Context) (id uuid.UUID, err error) {
	return uuid.Parse(c.Param("patient_id"))
}

func (v *PatientsController) IndexEndpoint(c *gin.Context) {
	patients, err := v.patientsRepository.All()
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientsResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, PatientsResponse{Patients: patients})
}

func (v *PatientsController) GetEndpoint(c *gin.Context) {
	patientID, err := extractPatientID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, PatientResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	var patient *domain.Patient
	patient, err = v.patientsRepository.FindByID(patientID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	c.JSON(http.StatusOK, PatientResponse{Patient: patient})
}

func (v *PatientsController) CreateEndpoint(c *gin.Context) {
	var (
		data inputPatient
		err  error
	)

	if err = c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, PatientResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	patient := data.buildModel()
	if err = v.patientsRepository.Create(&patient); err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	var newPatient *domain.Patient
	newPatient, err = v.patientsRepository.FindByID(patient.ID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}
	c.JSON(http.StatusCreated, PatientResponse{Patient: newPatient})
}

func (v *PatientsController) DeleteEndpoint(c *gin.Context) {
	var patient *domain.Patient
	id, err := extractPatientID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, PatientResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	patient, err = v.patientsRepository.FindByID(id)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	err = v.patientsRepository.Delete(patient)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}
	c.Status(http.StatusNoContent)
}

func (v *PatientsController) UpdateEndpoint(c *gin.Context) {
	var input inputPatient
	id, err := extractPatientID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, PatientResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, PatientResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	patient, err := v.patientsRepository.FindByID(id)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	input.updateModel(patient)
	err = v.patientsRepository.Update(patient)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	patient, err = v.patientsRepository.FindByID(id)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, PatientResponse{Error: r})
		return
	}

	c.JSON(http.StatusOK, PatientResponse{Patient: patient})
}
