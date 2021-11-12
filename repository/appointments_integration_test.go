package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/testutils"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AppointmentsIntegrationTestSuite struct {
	testutils.IntegrationSuite
	appointmentsRepository Appointments
}

func (s *AppointmentsIntegrationTestSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
	s.appointmentsRepository = NewAppointments(s.IntegrationSuite.DB(), zap.NewExample())
}

func (s *AppointmentsIntegrationTestSuite) TearDownSuite() {
	s.IntegrationSuite.TearDownSuite()
}

func (s *AppointmentsIntegrationTestSuite) TestCreate() {
	tests := []struct {
		name        string
		id          uuid.UUID
		appointment *domain.Appointment
		wantErr     error
	}{
		{
			name: "Successful",
			id:   uuid.MustParse("d93f7ecc-816f-4124-b41e-dcfa58f03761"),
			appointment: &domain.Appointment{
				ID:                uuid.MustParse("d93f7ecc-816f-4124-b41e-dcfa58f03761"),
				TreatmentCenterID: uuid.MustParse("52b2edf2-a380-4436-9f98-b70f78f174ef"),
				StartTime:         time.Date(2021, 11, 11, 8, 0, 0, 0, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "AlreadyExist",
			id:   uuid.MustParse("eecce415-2d4c-440d-ac90-9780a3bd3371"),
			appointment: &domain.Appointment{
				ID:                uuid.MustParse("eecce415-2d4c-440d-ac90-9780a3bd3371"),
				TreatmentCenterID: uuid.MustParse("52b2edf2-a380-4436-9f98-b70f78f174ef"),
				StartTime:         time.Date(2021, 11, 11, 8, 0, 0, 0, time.UTC),
			},
			wantErr: ErrUniqueConstraintFailure,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			err := s.appointmentsRepository.Create(tc.appointment)
			if tc.wantErr != nil {
				s.Assert().Equal(tc.wantErr, err)
			} else {
				gotAppointment, err := s.appointmentsRepository.FindByID(tc.id)
				s.Assert().NoError(err)

				s.Assert().Equal(tc.wantErr, err)
				s.Assert().Equal(tc.appointment.ID, gotAppointment.ID)
				s.Assert().Equal(tc.appointment.TreatmentCenterID, gotAppointment.TreatmentCenterID)
			}
		})
	}
}

func (s *AppointmentsIntegrationTestSuite) TestAllAvailable() {
	r, err := s.appointmentsRepository.AllAvailable()
	s.Assert().NoError(err)
	s.Assert().Len(r, 7)
}

func (s *AppointmentsIntegrationTestSuite) TestAllBookedByTreatmentCenterForDate() {
	r, err := s.appointmentsRepository.AllBookedByTreatmentCenterIDForDate(
		uuid.MustParse("10063726-d378-472c-9b50-22a48331635d"),
		time.Date(2021, 11, 13, 0, 0, 0, 0, time.UTC))
	s.Assert().NoError(err)
	s.Assert().Len(r, 1)
	s.Assert().Equal(uuid.MustParse("4cdb532d-bfe8-4af6-b9b5-d5078985a350"), r[0].ID)
	s.Assert().Equal(uuid.MustParse("10063726-d378-472c-9b50-22a48331635d"), r[0].TreatmentCenterID)
}

func TestAppointmentsIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping AppointmentsIntegrationTest in short mode.")
		return
	}
	t.Parallel()
	suite.Run(t, new(AppointmentsIntegrationTestSuite))
}
