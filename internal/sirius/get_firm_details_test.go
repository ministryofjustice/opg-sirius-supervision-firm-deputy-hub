package sirius

import (
	"bytes"
	"io/ioutil"
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
	[]
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
	"firmNumber": 100005
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := FirmDetails{
		ID: 2,
		FirmName: "Good Firm Inc",
		FirmNumber: 100005,
		Email: "good@firm.com",
		PhoneNumber: "123123123",
		AddressLine1: "10 St Hope Street",
		AddressLine2: "Wellington",
		AddressLine3: "",
		Town: "London",
		County: "Buckinghamshire",
		Postcode: "BU1 1TF",
	}

	firmDetails, err := client.GetFirmDetails(getContext(nil), 2)

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Equal(t, nil, err)
}

func TestGetFirmReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firmDetails, err := client.GetFirmDetails(getContext(nil), 1)

	expectedResponse := FirmDetails{
		ID: 0,
		FirmName: "",
		FirmNumber: 0,
		Email: "",
		PhoneNumber: "",
		AddressLine1: "",
		AddressLine2: "",
		AddressLine3: "",
		Town: "",
		County: "",
		Postcode: "",
	}

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/firm/1",
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

	expectedResponse := FirmDetails{
		ID: 0,
		FirmName: "",
		FirmNumber: 0,
		Email: "",
		PhoneNumber: "",
		AddressLine1: "",
		AddressLine2: "",
		AddressLine3: "",
		Town: "",
		County: "",
		Postcode: "",
	}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, firmDetails)
}
