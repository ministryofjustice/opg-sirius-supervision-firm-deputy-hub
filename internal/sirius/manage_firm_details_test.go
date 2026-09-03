package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"

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

	firmDetails := model.FirmDetails{
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

	err := client.ManageFirmDetails(getContext(nil), model.FirmDetails{ID: 1})

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

	err := client.ManageFirmDetails(getContext(nil), model.FirmDetails{ID: 1})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/1",
		Method: http.MethodPut,
	}, err)
}

func TestManageFirmReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.ManageFirmDetails(getContext(nil), model.FirmDetails{})

	assert.Equal(t, ErrUnauthorized, err)
}

func TestManageFirmDetails_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		UponReceiving("A request to edit firm details").
		WithRequest(http.MethodPut, SupervisionAPIPath+"/v1/firms/1", func(b *consumer.V2RequestBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":                   matchers.Like(1),
				"firmName":             matchers.Like("good firm inc"),
				"firmNumber":           matchers.Like(0),
				"email":                matchers.Like("good@firm.com"),
				"phoneNumber":          matchers.Like("077895526543"),
				"addressLine1":         matchers.Like("10 new street"),
				"addressLine2":         matchers.Like("new firm road"),
				"addressLine3":         matchers.Like("firmly"),
				"town":                 matchers.Like("Birmingham"),
				"county":               matchers.Like("Worcestershire"),
				"postcode":             matchers.Like("B1 1TF"),
				"executiveCaseManager": matchers.StructMatcher{"id": matchers.Like(0), "displayName": matchers.Like("")},
				"deputies":             matchers.Like([]model.FirmDeputies(nil)),
			})
		}).
		WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":           matchers.Like(1),
				"firmName":     matchers.Like("good firm inc"),
				"firmNumber":   matchers.Like(1000001),
				"email":        matchers.Like("good@firm.com"),
				"phoneNumber":  matchers.Like("077895526543"),
				"addressLine1": matchers.Like("10 new street"),
				"addressLine2": matchers.Like("new firm road"),
				"addressLine3": matchers.Like("firmly"),
				"town":         matchers.Like("Birmingham"),
				"county":       matchers.Like("Worcestershire"),
				"postcode":     matchers.Like("B1 1TF"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			return client.ManageFirmDetails(getContext(nil), model.FirmDetails{
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
			})
		})

	assert.NoError(t, err)
}
