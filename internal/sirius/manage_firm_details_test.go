package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestManageFirmDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"id":1,
		"firmName":"good firm inc",
		"firmNumber":1000001,
		"email":"good@firm.com",
		"phoneNumber":"077895526543",
		"addressLine1":"10 new street",
		"addressLine2":"new firm road",
		"addressLine3":"firmly",
		"town":"Birmingham",
		"county":"Worcestershire",
		"postcode":"B1 1TF"
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	firmDetails := FirmDetails{
		ID:           1,
		FirmName:     "good firm inc",
		Email:        "good@firm.com",
		PhoneNumber:  "077895526543",
		AddressLine1: "10 new street",
		AddressLine2: "new firm road",
		AddressLine3: "firmly",
		Town:         "Birmingham",
		County:       "Worcestershire",
		Postcode:     "B1 1TF",
	}

	err := client.ManageFirmDetails(getContext(nil), firmDetails)
	assert.Nil(t, err)
}

func TestManageFirmReturnsValidationError(t *testing.T) {
	client, _ := NewClient(&mocks.MockClient{}, "http://localhost:3000")

	json := `{"validation_errors": {"Test": {"error": "message"}}}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	err := client.ManageFirmDetails(getContext(nil), FirmDetails{ID: 1})

	assert.Equal(t, ValidationError{
		Errors: ValidationErrors{"Test": {"error": "message"}},
	}, err)
}

func TestManageFirmReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.ManageFirmDetails(getContext(nil), FirmDetails{ID: 1})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/firms/1",
		Method: http.MethodPut,
	}, err)
}

func TestManageFirmReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.ManageFirmDetails(getContext(nil), FirmDetails{})

	assert.Equal(t, ErrUnauthorized, err)
}
