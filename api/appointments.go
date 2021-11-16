package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/repository"
	"go.uber.org/zap"
)

type AppointmentsResponse struct {
	Appointments []*domain.Appointment `json:"appointments,omitempty"`
	Error        *ErrorResponse        `json:"error,omitempty"`
}

type AppointmentResponse struct {
	Appointment *domain.Appointment `json:"appointment,omitempty"`
	Error       *ErrorResponse      `json:"error,omitempty"`
}

type AppointmentBookingResponse struct {
	AppointmentBooking *domain.AppointmentBooking `json:"appointment_booking,omitempty"`
	Error              *ErrorResponse             `json:"error,omitempty"`
}

type AppointmentsController struct {
	appointmentsRepository        repository.Appointments
	appointmentBookingsRepository repository.AppointmentBookings
	logger                        *zap.Logger
}

type inputAppointment struct {
	domain.Appointment
}

func (t *inputAppointment) buildDomain() domain.Appointment {
	return domain.Appointment{
		ID:                uuid.New(),
		TreatmentCenterID: t.TreatmentCenterID,
		StartTime:         t.StartTime,
	}
}

func (t *inputAppointment) updateModel(appointment *domain.Appointment) {
	// appointment.Name = t.Name
}

type inputAppointmentBooking struct {
	domain.AppointmentBooking
}

func (t *inputAppointmentBooking) buildDomain(appointmentID uuid.UUID) domain.AppointmentBooking {
	return domain.AppointmentBooking{
		ID:            uuid.New(),
		PatientID:     t.PatientID,
		AppointmentID: appointmentID,
	}
}

func SetupAppointment(router gin.IRouter,
	appointmentsRepository repository.Appointments,
	appointmentBookingsRepository repository.AppointmentBookings,
	logger *zap.Logger) {
	c := AppointmentsController{
		appointmentsRepository:        appointmentsRepository,
		appointmentBookingsRepository: appointmentBookingsRepository,
		logger:                        logger.With(zap.String("component", "AppointmentsController")),
	}
	g := router.Group(
		"/appointments",
	)
	g.GET("/", c.IndexEndpoint)
	g.POST("/", c.CreateEndpoint)
	g.GET("/:appointment_id", c.GetEndpoint)
	g.POST("/:appointment_id/bookings", c.AddBookingEndpoint)
}

func extractAppointmentID(c *gin.Context) (id uuid.UUID, err error) {
	return uuid.Parse(c.Param("appointment_id"))
}

func (v *AppointmentsController) IndexEndpoint(c *gin.Context) {
	appointments, err := v.appointmentsRepository.AllAvailable()
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, AppointmentsResponse{Appointments: appointments})
}

func (v *AppointmentsController) GetEndpoint(c *gin.Context) {
	appointmentID, err := extractAppointmentID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, AppointmentResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	appointment, err := v.appointmentsRepository.FindByID(appointmentID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentResponse{Error: r})
		return
	}
	c.JSON(http.StatusOK, AppointmentResponse{Appointment: appointment})
}

func (v *AppointmentsController) CreateEndpoint(c *gin.Context) {
	var (
		input inputAppointment
		err   error
	)
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, AppointmentResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}
	appointment := input.buildDomain()
	err = v.appointmentsRepository.Create(&appointment)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentResponse{Error: r})
		return
	}
	var t *domain.Appointment
	t, err = v.appointmentsRepository.FindByID(appointment.ID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentResponse{Error: r})
		return
	}
	c.JSON(http.StatusCreated, AppointmentResponse{Appointment: t})
}

func (v *AppointmentsController) AddBookingEndpoint(c *gin.Context) {
	appointmentID, err := extractAppointmentID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, AppointmentBookingResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	var input inputAppointmentBooking
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, AppointmentBookingResponse{Error: &ErrorResponse{Msg: err.Error()}})
		return
	}

	appointmentBooking := input.buildDomain(appointmentID)
	appointmentBooking.Status = domain.AwaitingConfirmation
	err = v.appointmentBookingsRepository.Create(&appointmentBooking)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentBookingResponse{Error: r})
		return
	}

	var ab *domain.AppointmentBooking
	ab, err = v.appointmentBookingsRepository.FindByID(appointmentBooking.ID)
	if err != nil {
		code, r := handleRepositoryError(err, v.logger)
		c.JSON(code, AppointmentBookingResponse{Error: r})
		return
	}
	c.JSON(http.StatusCreated, AppointmentBookingResponse{AppointmentBooking: ab})
}
