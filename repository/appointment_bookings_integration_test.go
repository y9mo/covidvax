package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/testutils"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AppointmentBookingsIntegrationTestSuite struct {
	testutils.IntegrationSuite
	appointmentsRepository AppointmentBookings
}

func (s *AppointmentBookingsIntegrationTestSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
	s.appointmentsRepository = NewAppointmentBookings(s.IntegrationSuite.DB(), zap.NewExample())
}

func (s *AppointmentBookingsIntegrationTestSuite) TearDownSuite() {
	s.IntegrationSuite.TearDownSuite()
}

func (s *AppointmentBookingsIntegrationTestSuite) TestCreate() {
	tests := []struct {
		name        string
		id          uuid.UUID
		appointment *domain.AppointmentBooking
		wantErr     error
	}{
		{
			name: "Successful",
			id:   uuid.MustParse("755deca6-e171-4478-abfc-fa794830c7e1"),
			appointment: &domain.AppointmentBooking{
				ID:            uuid.MustParse("755deca6-e171-4478-abfc-fa794830c7e1"),
				AppointmentID: uuid.MustParse("eecce415-2d4c-440d-ac90-9780a3bd3371"),
				PatientID:     uuid.MustParse("8152fcbe-3228-46c9-b483-edcb6317d99c"),
				Status:        domain.Confirmed,
			},
			wantErr: nil,
		},
		{
			name: "IdAlreadyExist",
			id:   uuid.MustParse("f859ae2c-e24f-46e8-9c27-4431112fc710"),
			appointment: &domain.AppointmentBooking{
				ID:            uuid.MustParse("f859ae2c-e24f-46e8-9c27-4431112fc710"),
				AppointmentID: uuid.MustParse("dcea3eae-3004-4d3d-95b6-abc02ecb026d"),
				PatientID:     uuid.MustParse("8152fcbe-3228-46c9-b483-edcb6317d99c"),
				Status:        domain.AwaitingConfirmation,
			},
			wantErr: ErrUniqueConstraintFailure,
		},
		{
			name: "BookingAlreadyExist",
			id:   uuid.MustParse("0636fd0c-c524-4e42-b024-8824773a1c76"),
			appointment: &domain.AppointmentBooking{
				ID:            uuid.MustParse("0636fd0c-c524-4e42-b024-8824773a1c76"),
				AppointmentID: uuid.MustParse("4cdb532d-bfe8-4af6-b9b5-d5078985a350"),
				PatientID:     uuid.MustParse("24e32685-0a32-4a9d-bc22-0e98cdaf5884"),
				Status:        domain.AwaitingConfirmation,
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
				gotAppointmentBooking, err := s.appointmentsRepository.FindByID(tc.id)
				s.Assert().NoError(err)

				s.Assert().Equal(tc.wantErr, err)
				s.Assert().Equal(tc.appointment.ID, gotAppointmentBooking.ID)
				s.Assert().Equal(tc.appointment.AppointmentID, gotAppointmentBooking.AppointmentID)
				s.Assert().Equal(tc.appointment.PatientID, gotAppointmentBooking.PatientID)
			}
		})
	}
}

func TestAppointmentBookingsIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping AppointmentBookingsIntegrationTest in short mode.")
		return
	}
	t.Parallel()
	suite.Run(t, new(AppointmentBookingsIntegrationTestSuite))
}
