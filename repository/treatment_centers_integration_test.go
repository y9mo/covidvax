package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/testutils"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type TreatmentCentersIntegrationTestSuite struct {
	testutils.IntegrationSuite
	treatmentCentersRepository TreatmentCenters
}

func (s *TreatmentCentersIntegrationTestSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
	s.treatmentCentersRepository = NewTreatmentCenters(s.IntegrationSuite.DB(), zap.NewExample())
}

func (s *TreatmentCentersIntegrationTestSuite) TearDownSuite() {
	s.IntegrationSuite.TearDownSuite()
}

func (s *TreatmentCentersIntegrationTestSuite) TestCreate() {
	tests := []struct {
		name    string
		id      uuid.UUID
		patient *domain.TreatmentCenter
		wantErr error
	}{
		{
			name: "Successful",
			id:   uuid.MustParse("1e79b4d8-cae6-418a-bc95-0b6799540f3b"),
			patient: &domain.TreatmentCenter{
				ID:      uuid.MustParse("1e79b4d8-cae6-418a-bc95-0b6799540f3b"),
				Name:    "Center One",
				Address: "One place middle earth",
				Phone:   "0933420011",
			},
			wantErr: nil,
		},
		{
			name: "AlreadyExist",
			id:   uuid.MustParse("52b2edf2-a380-4436-9f98-b70f78f174ef"),
			patient: &domain.TreatmentCenter{
				ID:      uuid.MustParse("52b2edf2-a380-4436-9f98-b70f78f174ef"),
				Name:    "Center Two",
				Address: "Two is the right number",
				Phone:   "0422420033",
			},
			wantErr: ErrUniqueConstraintFailure,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			err := s.treatmentCentersRepository.Create(tc.patient)
			if tc.wantErr != nil {
				s.Assert().Equal(tc.wantErr, err)
			} else {
				gotTreatmentCenter, err := s.treatmentCentersRepository.FindByID(tc.id)
				s.Assert().NoError(err)

				s.Assert().Equal(tc.wantErr, err)
				s.Assert().Equal(tc.patient.ID, gotTreatmentCenter.ID)
				s.Assert().Equal(tc.patient.Name, gotTreatmentCenter.Name)
				s.Assert().Equal(tc.patient.Address, gotTreatmentCenter.Address)
				s.Assert().Equal(tc.patient.Phone, gotTreatmentCenter.Phone)
			}
		})
	}
}

func TestTreatmentCentersIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TreatmentCentersIntegrationTest in short mode.")
		return
	}
	t.Parallel()
	suite.Run(t, new(TreatmentCentersIntegrationTestSuite))
}
