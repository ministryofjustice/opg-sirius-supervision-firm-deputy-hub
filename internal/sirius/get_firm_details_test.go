package sirius

import (
	"bytes"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
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

func TestGetFirmDetails_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		UponReceiving("A request to get firm details").
		WithRequest(http.MethodGet, SupervisionAPIPath+"/v1/firms/2").
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"id":           matchers.Like(2),
				"firmName":     matchers.Like("Good Firm Inc"),
				"firmNumber":   matchers.Like(100005),
				"email":        matchers.Like("good@firm.com"),
				"phoneNumber":  matchers.Like("123123123"),
				"addressLine1": matchers.Like("10 St Hope Street"),
				"addressLine2": matchers.Like("Wellington"),
				"addressLine3": matchers.Like(""),
				"town":         matchers.Like("London"),
				"county":       matchers.Like("Buckinghamshire"),
				"postcode":     matchers.Like("BU1 1TF"),
				"deputies": matchers.EachLike(matchers.StructMatcher{
					"id":               matchers.Like(77),
					"deputyNumber":     matchers.Like(22),
					"organisationName": matchers.Like("pro dept"),
				}, 1),
				"executiveCaseManager": matchers.StructMatcher{
					"id":          matchers.Like(71),
					"displayName": matchers.Like("LayTeam1 User1"),
				},
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			_, err := client.GetFirmDetails(getContext(nil), 2)
			return err
		})

	assert.NoError(t, err)
}
