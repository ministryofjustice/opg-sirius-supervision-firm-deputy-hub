package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFirmDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `	{
		"id": 2,
		"deputies": [
			{
				"id":77,
				"personType":"Deputy",
				"deputyStatus":"Inactive",
				"deputyNumber":22,
				"organisationName":"pro dept",
				"deputySubType":[]
			},
			{
				"id":75,
				"personType":"Deputy",
				"deputyStatus":"Active",
				"eveningNumber":"07748933233",
				"deputyNumber":20,
				"organisationName":"deputy pro",
				"organisationTeamOrDepartmentName":""
			}
		],
		"firmName": "Good Firm Inc",
		"addressLine1": "10 St Hope Street",
		"addressLine2": "Wellington",
		"addressLine3": "",
		"town": "London",
		"county": "Buckinghamshire",
		"postcode": "BU1 1TF",
		"phoneNumber": "123123123",
		"email": "good@firm.com",
		"firmNumber": 100005,
		"executiveCaseManager": {
			"id": 71,
			"name": "LayTeam1",
			"displayName": "LayTeam1 User1",
			"surname": "User1"
		}

	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.FirmDetails{
		ID:           2,
		FirmName:     "Good Firm Inc",
		FirmNumber:   100005,
		Email:        "good@firm.com",
		PhoneNumber:  "123123123",
		AddressLine1: "10 St Hope Street",
		AddressLine2: "Wellington",
		AddressLine3: "",
		Town:         "London",
		County:       "Buckinghamshire",
		Postcode:     "BU1 1TF",
		Deputies: []model.FirmDeputies{
			{
				DeputyId:         77,
				DeputyNumber:     22,
				OrganisationName: "pro dept",
			},
			{
				DeputyId:         75,
				DeputyNumber:     20,
				OrganisationName: "deputy pro",
			},
		},
		ExecutiveCaseManager: model.ExecutiveCaseManager{
			Id:          71,
			DisplayName: "LayTeam1 User1",
		},
	}

	firmDetails, err := client.GetFirmDetails(getContext(nil), 2)

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Nil(t, err)
}

func TestGetFirmReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firmDetails, err := client.GetFirmDetails(getContext(nil), 1)

	expectedResponse := model.FirmDetails{}

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/1",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firmDetails, err := client.GetFirmDetails(getContext(nil), 1)

	expectedResponse := model.FirmDetails{}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, firmDetails)
}
