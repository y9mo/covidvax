package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/repository"
	"go.uber.org/zap"
)

type TreatmentCentersResponse struct {
	TreatmentCenters []*domain.TreatmentCenter `json:"treatment_centers,omitempty"`
	Error            *ErrorResponse            `json:"error,omitempty"`
}

type TreatmentCenterResponse struct {
	TreatmentCenter *domain.TreatmentCenter `json:"treatment_center,omitempty"`
	Error           *ErrorResponse          `json:"error,omitempty"`
}

type TreatmentCenterAppointmentResponse struct {
	Appointments *domain.Appointment `json:"appointments,omitempty"`
	Error        *ErrorResponse      `json:"error,omitempty"`
}

type TreatmentCenterAppointmentsResponse struct {
	Appointments []*domain.Appointment `json:"appointments,omitempty"`
	Error        *ErrorResponse        `json:"error,omitempty"`
}

type TreatmentCentersController struct {
	treatmentCentersRepository repository.TreatmentCenters
	appointmentsRepository     repository.Appointments
	logger                     *zap.Logger
}

type inputTreatmentCenter struct {
	domain.TreatmentCenter
}

func (t *inputTreatmentCenter) buildModel() domain.TreatmentCenter {
	return domain.TreatmentCenter{
		ID:   uuid.New(),
		Name: t.Name,
	}
}

func (t *inputTreatmentCenter) updateModel(treatmentCenter *domain.TreatmentCenter) {
	treatmentCenter.Name = t.Name
	treatmentCenter.Address = t.Address
}

type TreatmentCenterAppointmentRequest struct {
	Date *time.Time `form:"date" time_format:"2006-01-02" time_utc:"1"`
}

func SetupTreatmentCenter(
	router gin.IRouter,
	treatmentCentersRepository repository.TreatmentCenters,
	appointmentsRepository repository.Appointments,
	logger *zap.Logger) {
	c := TreatmentCentersController{
		treatmentCentersRepository: treatmentCentersRepository,
		appointmentsRepository:     appointmentsRepository,
		logger:                     logger.With(zap.String("component", "TreatmentCentersController")),
	}
	g := router.Group(
		"/treatment_centers",
	)
	g.GET("/", c.IndexEndpoint)
	g.POST("/", c.CreateEndpoint)
	g.GET("/:treatment_center_id", c.GetEndpoint)
	g.GET("/:treatment_center_id/bookings", c.GetBookedAppointmentsEndpoint)
}

func extractTreatmentCenterID(c *gin.Context) (id uuid.UUID, err error) {
	return uuid.Parse(c.Param("treatment_center_id"))
}

func currentDaytime() *time.Time {
	d := time.Now().UTC()
	return &d
}

func (v *TreatmentCentersController) IndexEndpoint(c *gin.Context) {
	treatmentCenters, err := v.treatmentCentersRepository.All()
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, TreatmentCenterResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, TreatmentCentersResponse{TreatmentCenters: treatmentCenters})
}

func (v *TreatmentCentersController) GetEndpoint(c *gin.Context) {
	id, err := extractTreatmentCenterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, TreatmentCenterResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	treatmentCenter, err := v.treatmentCentersRepository.FindByID(id)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, TreatmentCenterResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, TreatmentCenterResponse{TreatmentCenter: treatmentCenter})
}

func (v *TreatmentCentersController) CreateEndpoint(c *gin.Context) {
	var (
		input inputTreatmentCenter
		err   error
	)

	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, TreatmentCenterResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	treatmentCenter := input.buildModel()
	err = v.treatmentCentersRepository.Create(&treatmentCenter)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, TreatmentCenterResponse{Error: r})
		return
	}
	var t *domain.TreatmentCenter
	t, err = v.treatmentCentersRepository.FindByID(treatmentCenter.ID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, TreatmentCenterResponse{Error: r})
		return
	}
	c.JSON(http.StatusCreated, TreatmentCenterResponse{TreatmentCenter: t})
}

func (v *TreatmentCentersController) GetBookedAppointmentsEndpoint(c *gin.Context) {
	id, err := extractTreatmentCenterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, TreatmentCenterAppointmentsResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	var request TreatmentCenterAppointmentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, TreatmentCenterAppointmentsResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	if request.Date == nil {
		request.Date = currentDaytime()
	}

	treatmentCenterAppointments, err := v.appointmentsRepository.AllBookedByTreatmentCenterIDForDate(id,
		*request.Date)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, TreatmentCenterAppointmentsResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, TreatmentCenterAppointmentsResponse{Appointments: treatmentCenterAppointments})
}
