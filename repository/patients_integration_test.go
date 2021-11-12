package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/y9mo/covidvax/domain"
	"github.com/y9mo/covidvax/testutils"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type PatientsIntegrationTestSuite struct {
	testutils.IntegrationSuite
	patientsRepository Patients
}

func (s *PatientsIntegrationTestSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
	s.patientsRepository = NewPatients(s.IntegrationSuite.DB(), zap.NewExample())
}

func (s *PatientsIntegrationTestSuite) TearDownSuite() {
	s.IntegrationSuite.TearDownSuite()
}

func (s *PatientsIntegrationTestSuite) TestCreate() {
	tests := []struct {
		name    string
		id      uuid.UUID
		patient *domain.Patient
		wantErr error
	}{
		{
			name: "Successful",
			id:   uuid.MustParse("d93f7ecc-816f-4124-b41e-dcfa58f03761"),
			patient: &domain.Patient{
				ID:        uuid.MustParse("d93f7ecc-816f-4124-b41e-dcfa58f03761"),
				Email:     "patient.zero@some.com",
				FirstName: "Patient",
				LastName:  "Zero",
			},
			wantErr: nil,
		},
		{
			name: "AlreadyExist",
			id:   uuid.MustParse("cf1d40a0-055d-4116-a66f-9bb1624e66fd"),
			patient: &domain.Patient{
				ID:        uuid.MustParse("cf1d40a0-055d-4116-a66f-9bb1624e66fd"),
				Email:     "patient.one@some.com",
				FirstName: "Patient",
				LastName:  "One",
			},
			wantErr: ErrUniqueConstraintFailure,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			err := s.patientsRepository.Create(tc.patient)
			if tc.wantErr != nil {
				s.Assert().Equal(tc.wantErr, err)
			} else {
				gotPatient, err := s.patientsRepository.FindByID(tc.id)
				s.Assert().NoError(err)

				s.Assert().Equal(tc.wantErr, err)
				s.Assert().Equal(tc.patient.ID, gotPatient.ID)
				s.Assert().Equal(tc.patient.Email, gotPatient.Email)
				s.Assert().Equal(tc.patient.FirstName, gotPatient.FirstName)
				s.Assert().Equal(tc.patient.LastName, gotPatient.LastName)
			}
		})
	}
}

func TestPatientsIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping PatientsIntegrationTest in short mode.")
		return
	}
	t.Parallel()
	suite.Run(t, new(PatientsIntegrationTestSuite))
}
