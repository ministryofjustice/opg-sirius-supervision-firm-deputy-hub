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

func TestEditPii(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"piiReceived":"20/01/2020",
		"piiExpiry":"20/01/2025",
		"piiAmount":254,
		"piiRequested":"10/01/2020"
		}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	piiDetails := model.PiiDetails{
		FirmId:       21,
		PiiReceived:  "20/01/2020",
		PiiExpiry:    "20/01/2025",
		PiiAmount:    254,
		PiiRequested: "10/01/2020",
	}

	err := client.EditPiiCertificate(getContext(nil), piiDetails)
	assert.Nil(t, err)
}

func TestEditPiiReturnsValidationError(t *testing.T) {
	client, _ := NewClient(&mocks.MockClient{}, "http://localhost:3000")

	json := `{"validation_errors": {"Test": {"error": "message"}}}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	err := client.EditPiiCertificate(getContext(nil), model.PiiDetails{})

	assert.Equal(t, ValidationError{
		Errors: ValidationErrors{"Test": {"error": "message"}},
	}, err)
}

func TestEditPiiReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.EditPiiCertificate(getContext(nil), model.PiiDetails{})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/0/indemnity-insurance",
		Method: http.MethodPut,
	}, err)
}

func TestEditPiiReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.EditPiiCertificate(getContext(nil), model.PiiDetails{})

	assert.Equal(t, ErrUnauthorized, err)

}

func TestEditPii_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("User exists").
		UponReceiving("A request to edit PII").
		WithRequest(http.MethodPut, SupervisionAPIPath+"/v1/firms/21/indemnity-insurance", func(b *consumer.V2RequestBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"firmId":       matchers.Like(21),
				"piiReceived":  matchers.Like("20/01/2020"),
				"piiExpiry":    matchers.Like("20/01/2025"),
				"piiAmount":    matchers.Like(254),
				"piiRequested": matchers.Like("10/01/2020"),
			})
		}).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"firmId":       matchers.Like(21),
				"piiReceived":  matchers.Like("20/01/2020"),
				"piiExpiry":    matchers.Like("20/01/2025"),
				"piiAmount":    matchers.Like(254),
				"piiRequested": matchers.Like("10/01/2020"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			return client.EditPiiCertificate(getContext(nil), model.PiiDetails{
				FirmId:       21,
				PiiReceived:  "20/01/2020",
				PiiExpiry:    "20/01/2025",
				PiiAmount:    254,
				PiiRequested: "10/01/2020",
			})
		})

	assert.NoError(t, err)
}
