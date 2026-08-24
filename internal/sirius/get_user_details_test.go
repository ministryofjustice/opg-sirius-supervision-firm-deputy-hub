package sirius

import (
	"bytes"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetUserDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
	"id": 68,
	"roles": ["Finance Manager", "System Admin"]
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Assignee{
		ID:    68,
		Roles: []string{"Finance Manager", "System Admin"},
	}

	userDetails, err := client.GetUserDetails(getContext(nil))

	assert.Equal(t, expectedResponse, userDetails)
	assert.Equal(t, nil, err)
}

func TestUserDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	userDetails, err := client.GetUserDetails(getContext(nil))

	expectedResponse := model.Assignee{ID: 0}

	assert.Equal(t, expectedResponse, userDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestUserDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	userDetails, err := client.GetUserDetails(getContext(nil))

	expectedResponse := model.Assignee{ID: 0}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, userDetails)
}

func TestGetUserDetails_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)
	err = pact.
		AddInteraction().
		Given("I am a System Admin").
		UponReceiving("A request to get the current user details").
		WithRequest(http.MethodGet, SupervisionAPIPath+"/v1/users/current").
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":          matchers.Like(123),
				"displayName": matchers.Like("Ian Finance"),
				"roles":       matchers.EachLike("Finance Manager", 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))

			userDetails, err := client.GetUserDetails(getContext(nil))
			if err != nil {
				return err
			}

			assert.Equal(t, "Finance Manager", userDetails.Roles[0])
			return nil
		})
	assert.NoError(t, err)
}
