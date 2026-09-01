package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestChangeECM(t *testing.T) {
	client, _ := NewClient(&mocks.MockClient{}, "http://localhost:3000")

	json := `{"ecmId": 32}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, model.FirmDetails{ID: 76})
	assert.Equal(t, nil, err)
}

func TestChangeECMReturnsValidationError(t *testing.T) {
	client, _ := NewClient(&mocks.MockClient{}, "http://localhost:3000")

	json := `{"validation_errors": {"Test": {"error": "message"}}}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}

	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, model.FirmDetails{ID: 0})

	assert.Equal(t, ValidationError{
		Errors: ValidationErrors{"Test": {"error": "message"}},
	}, err)
}

func TestChangeECMReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, model.FirmDetails{ID: 76})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/76/ecm",
		Method: http.MethodPut,
	}, err)
}

func TestChangeECMReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, model.FirmDetails{ID: 76})

	assert.Equal(t, ErrUnauthorized, err)
}

func TestChangeECM_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		UponReceiving("A request to change a firms ECM").
		WithRequest(http.MethodPut, SupervisionAPIPath+"/v1/firms/76/ecm", func(b *consumer.V2RequestBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"ecmId": matchers.Like(23),
			})
		}).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":       matchers.Like(76),
				"firmName": matchers.Like("Example Firm"),
				"executiveCaseManager": matchers.Like(map[string]any{
					"id": matchers.Like(32),
				}),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			return client.ChangeECM(getContext(nil), ExecutiveCaseManagerOutgoing{EcmId: 23}, model.FirmDetails{ID: 76})
		})

	assert.NoError(t, err)
}
