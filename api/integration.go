package api

import (
	"github.com/gin-gonic/gin"
	"github.com/y9mo/covidvax/repository"
	"github.com/y9mo/covidvax/testutils"
	"go.uber.org/zap"
)

type ApiIntegrationSuite struct {
	testutils.IntegrationSuite
	Router *gin.Engine
}

func (s *ApiIntegrationSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
	var err error

	logger := zap.NewExample()

	pr := repository.NewPatients(s.DB(), logger)
	tcr := repository.NewTreatmentCenters(s.DB(), logger)
	ar := repository.NewAppointments(s.DB(), logger)
	abr := repository.NewAppointmentBookings(s.DB(), logger)

	s.Router, err = Setup(logger, pr, tcr, ar, abr)
	s.Require().NoError(err)
}
