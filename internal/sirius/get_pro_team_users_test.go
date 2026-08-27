package sirius

import (
	"bytes"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPaDeputyTeamUsersReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
	{
		"id": 25,
		"name": "Pro Team 1 - (Supervision)",
		"phoneNumber": "0123456789",
		"displayName": "Pro Team 1 - (Supervision)",
		"deleted": false,
		"email": "ProTeam1.team@opgtest.com",
		"members": [
			{
				"id": 90,
				"name": "LayTeam1",
				"phoneNumber": "12345678",
				"displayName": "LayTeam1 User20",
				"deleted": false,
				"email": "lay1-20@opgtest.com"
			},
			{
				"id": 94,
				"name": "PROTeam1",
				"phoneNumber": "12345678",
				"displayName": "PROTeam1 User1",
				"deleted": false,
				"email": "pro1@opgtest.com"
			}
		],
		"groupName": null,
		"parent": null,
		"children": [],
		"teamType": {
			"handle": "PRO",
			"label": "Pro",
			"deprecated": null
		}
	},
	{
		"id": 26,
		"name": "Pro Team 2 - (Supervision)",
		"phoneNumber": "0123456789",
		"displayName": "Pro Team 2 - (Supervision)",
		"deleted": false,
		"email": "ProTeam2.team@opgtest.com",
		"members": [
			{
				"id": 37,
				"name": "atwo",
				"phoneNumber": "03004560300",
				"displayName": "atwo manager",
				"deleted": false,
				"email": "2manager@opgtest.com"
			},
			{
				"id": 101,
				"name": "CardPayment",
				"phoneNumber": "12345678",
				"displayName": "CardPayment User",
				"deleted": false,
				"email": "card.payment.user@opgtest.com"
			}
		],
		"groupName": null,
		"parent": null,
		"children": [],
		"teamType": {
			"handle": "PRO",
			"label": "Pro",
			"deprecated": null
		}
	}
]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.Member{
		{
			Id:          90,
			DisplayName: "LayTeam1 User20",
		},
		{
			Id:          94,
			DisplayName: "PROTeam1 User1",
		},
		{
			Id:          37,
			DisplayName: "atwo manager",
		},
		{
			Id:          101,
			DisplayName: "CardPayment User",
		},
	}

	_, proDeputyTeam, err := client.GetProTeamUsers(getContext(nil))

	assert.Equal(t, expectedResponse, proDeputyTeam)
	assert.Equal(t, nil, err)
}

func TestGetPaDeputyTeamUsersReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	_, proDeputyMembers, err := client.GetProTeamUsers(getContext(nil))

	var expectedResponse []model.Member

	assert.Equal(t, expectedResponse, proDeputyMembers)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/teams?type=pro",
		Method: http.MethodGet,
	}, err)
}

func TestGetPaDeputyTeamUsersReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	_, proDeputyMembers, err := client.GetProTeamUsers(getContext(nil))

	var expectedResponse []model.Member

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, proDeputyMembers)
}

func TestGetProTeamUsers_contract(t *testing.T) {
	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "supervision-firm-deputy-hub",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		UponReceiving("A request to get pro teams").
		WithRequest(http.MethodGet, SupervisionAPIPath+"/v1/teams", func(b *consumer.V2RequestBuilder) {
			b.Query("type", matchers.S("pro"))
		}).
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody([]interface{}{
				map[string]interface{}{
					"id":          25,
					"name":        "Pro Team 1 - (Supervision)",
					"displayName": "Pro Team 1 - (Supervision)",
					"phoneNumber": "0123456789",
					"deleted":     false,
					"email":       "ProTeam1.team@opgtest.com",
					"members": []interface{}{
						map[string]interface{}{
							"id":          90,
							"name":        "LayTeam1",
							"displayName": "LayTeam1 User20",
							"phoneNumber": "12345678",
							"deleted":     false,
							"email":       "lay1-20@opgtest.com",
						},
						map[string]interface{}{
							"id":          94,
							"name":        "PROTeam1",
							"displayName": "PROTeam1 User1",
							"phoneNumber": "12345678",
							"deleted":     false,
							"email":       "pro1@opgtest.com",
						},
					},
					"teamType": map[string]interface{}{
						"handle": "PRO",
						"label":  "Pro",
					},
				},
				map[string]interface{}{
					"id":          26,
					"name":        "Pro Team 2 - (Supervision)",
					"displayName": "Pro Team 2 - (Supervision)",
					"phoneNumber": "0123456789",
					"deleted":     false,
					"email":       "ProTeam2.team@opgtest.com",
					"members": []interface{}{
						map[string]interface{}{
							"id":          37,
							"name":        "atwo",
							"displayName": "atwo manager",
							"phoneNumber": "03004560300",
							"deleted":     false,
							"email":       "2manager@opgtest.com",
						},
						map[string]interface{}{
							"id":          101,
							"name":        "CardPayment",
							"displayName": "CardPayment User",
							"phoneNumber": "12345678",
							"deleted":     false,
							"email":       "card.payment.user@opgtest.com",
						},
					},
					"teamType": map[string]interface{}{
						"handle": "PRO",
						"label":  "Pro",
					},
				},
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			_, members, err := client.GetProTeamUsers(getContext(nil))
			if err != nil {
				return err
			}
			assert.Equal(t, []model.Member{
				{Id: 90, DisplayName: "LayTeam1 User20"},
				{Id: 94, DisplayName: "PROTeam1 User1"},
				{Id: 37, DisplayName: "atwo manager"},
				{Id: 101, DisplayName: "CardPayment User"},
			}, members)
			return nil
		})

	assert.NoError(t, err)
}
