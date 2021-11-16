package api

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/y9mo/covidvax"
	"github.com/y9mo/covidvax/repository"
	"go.uber.org/zap"
)

func Setup(
	logger *zap.Logger,
	pr repository.Patients,
	tcr repository.TreatmentCenters,
	ar repository.Appointments,
	abr repository.AppointmentBookings,
) (*gin.Engine, error) {

	router := gin.New()

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Logs panic to error log
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.GET("/", Index)

	g := router.Group("/v1")
	SetupPatient(g, pr, logger)
	SetupTreatmentCenter(g, tcr, ar, logger)
	SetupAppointment(g, ar, abr, logger)
	return router, nil
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"version":   covidvax.Version,
		"component": "covidvax",
	})
}
