package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestRequestPii(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"firmId":2,
		"piiRequested":"10/01/2020"
		}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	piiDetails := PiiDetailsRequest{
		FirmId:       2,
		PiiRequested: "10/01/2020",
	}

	err := client.RequestPiiCertificate(getContext(nil), piiDetails)
	assert.Nil(t, err)
}

func TestRequestPiiReturnsValidationError(t *testing.T) {
	client, _ := NewClient(&mocks.MockClient{}, "http://localhost:3000")

	json := `{"validation_errors": {"Test": {"error": "message"}}}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	err := client.RequestPiiCertificate(getContext(nil), PiiDetailsRequest{})

	assert.Equal(t, ValidationError{
		Errors: ValidationErrors{"Test": {"error": "message"}},
	}, err)
}

func TestRequestPiiReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.RequestPiiCertificate(getContext(nil), PiiDetailsRequest{})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/0/indemnity-insurance",
		Method: http.MethodPatch,
	}, err)
}

func TestRequestPiiReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	err := client.RequestPiiCertificate(getContext(nil), PiiDetailsRequest{})

	assert.Equal(t, ErrUnauthorized, err)

}

func TestRequestPii_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		UponReceiving("A request to patch PII").
		WithRequest(http.MethodPatch, SupervisionAPIPath+"/v1/firms/2/indemnity-insurance", func(b *consumer.V2RequestBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"firmId":       matchers.Like(2),
				"piiRequested": matchers.Like("10/01/2020"),
			})
		}).
		WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"firmId":       matchers.Like(2),
				"piiRequested": matchers.Like("10/01/2020"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			return client.RequestPiiCertificate(getContext(nil), PiiDetailsRequest{
				FirmId:       2,
				PiiRequested: "10/01/2020",
			})
		})

	assert.NoError(t, err)
}
