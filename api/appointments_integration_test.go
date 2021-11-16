package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/suite"
	"github.com/y9mo/covidvax/testutils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type AppointmentsApiIntegrationTestSuite struct {
	ApiIntegrationSuite
}

const validAppointmentJSON = `{
	"treatment_center_id": "52b2edf2-a380-4436-9f98-b70f78f174ef",
	"start_time": "2021-11-14T10:00:00Z"
}`

const validAppointmentBookingJSON = `{
	"patient_id": "8152fcbe-3228-46c9-b483-edcb6317d99c"
}`

func (suite *AppointmentsApiIntegrationTestSuite) TestAppointmentsIndexList() {
	apitest.New().Debug().
		Handler(suite.Router).
		Get("/v1/appointments/").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.appointments`, 7)).
		Assert(jsonpath.Contains(`$.appointments[? @.id=="eecce415-2d4c-440d-ac90-9780a3bd3371"].treatment_center_id`,
			"52b2edf2-a380-4436-9f98-b70f78f174ef")).
		Assert(jsonpath.Contains(`$.appointments[? @.id=="5edac3af-6805-469e-ae94-f9610c09516a"].treatment_center_id`,
			"10063726-d378-472c-9b50-22a48331635d")).
		End()
}

func (suite *AppointmentsApiIntegrationTestSuite) TestCreateAppointment() {
	var id string
	apitest.New().Debug().
		Handler(suite.Router).
		Post("/v1/appointments/").
		JSON(validAppointmentJSON).
		Expect(suite.T()).
		Status(http.StatusCreated).
		Assert(jsonpath.Present(`$.appointment.id`)).
		Assert(testutils.Extract(`$.appointment.id`, &id)).
		Assert(jsonpath.Equal(`$.appointment.treatment_center_id`, "52b2edf2-a380-4436-9f98-b70f78f174ef")).
		Assert(jsonpath.Present(`$.appointment.created_at`)).
		Assert(jsonpath.Present(`$.appointment.updated_at`)).
		End()

	suite.Assert().NotEmpty(id)

	apitest.New().Debug().
		Handler(suite.Router).
		Get(fmt.Sprintf("/v1/appointments/%s", id)).
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.appointment.id`)).
		Assert(jsonpath.Equal(`$.appointment.id`, id)).
		Assert(jsonpath.Equal(`$.appointment.treatment_center_id`, "52b2edf2-a380-4436-9f98-b70f78f174ef")).
		Assert(jsonpath.Present(`$.appointment.created_at`)).
		Assert(jsonpath.Present(`$.appointment.updated_at`)).
		End()
}

func (suite *AppointmentsApiIntegrationTestSuite) TestGetAppointmentsByID() {
	apitest.New().Debug().
		Handler(suite.Router).
		Get("/v1/appointments/4cdb532d-bfe8-4af6-b9b5-d5078985a350").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.appointment.id`)).
		Assert(jsonpath.Equal(`$.appointment.id`, "4cdb532d-bfe8-4af6-b9b5-d5078985a350")).
		Assert(jsonpath.Equal(`$.appointment.treatment_center_id`, "10063726-d378-472c-9b50-22a48331635d")).
		Assert(jsonpath.Present(`$.appointment.created_at`)).
		Assert(jsonpath.Present(`$.appointment.updated_at`)).
		End()
}

func (suite *AppointmentsApiIntegrationTestSuite) TestCreateAppointmentBooking() {
	var id string
	apitest.New().Debug().
		Handler(suite.Router).
		Post("/v1/appointments/eecce415-2d4c-440d-ac90-9780a3bd3371/bookings").
		JSON(validAppointmentBookingJSON).
		Expect(suite.T()).
		Status(http.StatusCreated).
		Assert(jsonpath.Present(`$.appointment_booking.id`)).
		Assert(testutils.Extract(`$.appointment_booking.id`, &id)).
		Assert(jsonpath.Equal(`$.appointment_booking.appointment_id`, "eecce415-2d4c-440d-ac90-9780a3bd3371")).
		Assert(jsonpath.Equal(`$.appointment_booking.patient_id`, "8152fcbe-3228-46c9-b483-edcb6317d99c")).
		Assert(jsonpath.Equal(`$.appointment_booking.status`, "awaiting confirmation")).
		End()

	suite.Assert().NotEmpty(id)
}

func TestAppointmentsApiIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AppointmentsApiIntegrationTestSuite))
}
