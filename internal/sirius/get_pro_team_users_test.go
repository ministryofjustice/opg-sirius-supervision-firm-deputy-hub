package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPaDeputyTeamUsersReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
    "id": 23,
    "name": "PA Team 1 - (Supervision)",
    "phoneNumber": "0123456789",
    "displayName": "PA Team 1 - (Supervision)",
    "deleted": false,
    "email": "PATeam1.team@opgtest.com",
    "members": [
        {
            "id": 92,
            "name": "PATeam1",
            "phoneNumber": "12345678",
            "displayName": "PATeam1 User1",
            "deleted": false,
            "email": "pa1@opgtest.com"
        }
    ],
    "children": [],
    "teamType": {
        "handle": "PA",
        "label": "PA"
    }
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []TeamMember{
		{
			ID:          92,
			DisplayName: "PATeam1 User1",
		},
	}

	paDeputyTeam, err := client.GetPaDeputyTeamMembers(getContext(nil), 23)

	assert.Equal(t, expectedResponse, paDeputyTeam)
	assert.Equal(t, nil, err)
}

func TestGetPaDeputyTeamUsersReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	paDeputyTeam, err := client.GetPaDeputyTeamMembers(getContext(nil), 23)

	expectedResponse := []TeamMember([]TeamMember{})

	assert.Equal(t, expectedResponse, paDeputyTeam)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/teams/23",
		Method: http.MethodGet,
	}, err)
}

func TestGetPaDeputyTeamUsersReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	paDeputyTeam, err := client.GetPaDeputyTeamMembers(getContext(nil), 23)

	expectedResponse := []TeamMember([]TeamMember{})

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, paDeputyTeam)
}
