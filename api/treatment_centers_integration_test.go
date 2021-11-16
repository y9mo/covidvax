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

type TreatmentCentersApiIntegrationTestSuite struct {
	ApiIntegrationSuite
}

const validTreatmentCenterJSON = `{
	"id": "e46e6fff-3cc9-42cd-839d-e0528166b40b",
	"name": "Center for test",
	"address": "test",
	"phone": "0102030405"
}`

const validTreatmentCenterBookingJSON = `{
	"patient_id": "8152fcbe-3228-46c9-b483-edcb6317d99c"
}`

func (s *TreatmentCentersApiIntegrationTestSuite) TestTreatmentCentersIndexList() {
	apitest.New().Debug().
		Handler(s.Router).
		Get("/v1/treatment_centers/").
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.treatment_centers`, 4)).
		Assert(jsonpath.Contains(`$.treatment_centers[? @.id=="52b2edf2-a380-4436-9f98-b70f78f174ef"].name`,
			"Center Two")).
		Assert(jsonpath.Contains(`$.treatment_centers[? @.id=="52b2edf2-a380-4436-9f98-b70f78f174ef"].address`,
			"somewhere")).
		Assert(jsonpath.Contains(`$.treatment_centers[? @.id=="10063726-d378-472c-9b50-22a48331635d"].name`,
			"Center in the game")).
		Assert(jsonpath.Contains(`$.treatment_centers[? @.id=="10063726-d378-472c-9b50-22a48331635d"].address`,
			"game time")).
		End()
}

func (s *TreatmentCentersApiIntegrationTestSuite) TestCreateTreatmentCenter() {
	var id string
	apitest.New().Debug().
		Handler(s.Router).
		Post("/v1/treatment_centers/").
		JSON(validTreatmentCenterJSON).
		Expect(s.T()).
		Status(http.StatusCreated).
		Assert(jsonpath.Present(`$.treatment_center.id`)).
		Assert(testutils.Extract(`$.treatment_center.id`, &id)).
		Assert(jsonpath.Equal(`$.treatment_center.name`, "Center for test")).
		Assert(jsonpath.Present(`$.treatment_center.created_at`)).
		Assert(jsonpath.Present(`$.treatment_center.updated_at`)).
		End()

	s.Assert().NotEmpty(id)

	apitest.New().Debug().
		Handler(s.Router).
		Get(fmt.Sprintf("/v1/treatment_centers/%s", id)).
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.treatment_center.id`)).
		Assert(jsonpath.Equal(`$.treatment_center.id`, id)).
		Assert(jsonpath.Equal(`$.treatment_center.name`, "Center for test")).
		Assert(jsonpath.Present(`$.treatment_center.created_at`)).
		Assert(jsonpath.Present(`$.treatment_center.updated_at`)).
		End()
}

func (s *TreatmentCentersApiIntegrationTestSuite) TestGetTreatmentCentersByID() {
	apitest.New().Debug().
		Handler(s.Router).
		Get("/v1/treatment_centers/32b2edf2-a380-4436-9f98-b70f78f1934d").
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.treatment_center.id`)).
		Assert(jsonpath.Equal(`$.treatment_center.id`, "32b2edf2-a380-4436-9f98-b70f78f1934d")).
		Assert(jsonpath.Equal(`$.treatment_center.name`, "Center of Light")).
		Assert(jsonpath.Present(`$.treatment_center.created_at`)).
		Assert(jsonpath.Present(`$.treatment_center.updated_at`)).
		End()
}

func (s *TreatmentCentersApiIntegrationTestSuite) TestGetTreatmentCenterBookings() {
	apitest.New().Debug().
		Handler(s.Router).
		Get("/v1/treatment_centers/10063726-d378-472c-9b50-22a48331635d/bookings").
		QueryParams(map[string]string{"date": "2021-11-13"}).
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.appointments`)).
		Assert(jsonpath.Equal(`$.appointments[0].id`, "4cdb532d-bfe8-4af6-b9b5-d5078985a350")).
		Assert(jsonpath.Equal(`$.appointments[0].treatment_center_id`, "10063726-d378-472c-9b50-22a48331635d")).
		Assert(jsonpath.Present(`$.appointments[0].start_time`)).
		Assert(jsonpath.Present(`$.appointments[0].created_at`)).
		Assert(jsonpath.Present(`$.appointments[0].updated_at`)).
		End()
}

func TestTreatmentCentersApiIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TreatmentCentersApiIntegrationTestSuite))
}
