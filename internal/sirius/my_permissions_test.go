package sirius

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func teapotServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}),
	)
}

func getContext(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}

func TestUserPermissionsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	// build response JSON
	json := `{
 		"v1-users": {
 		  "permissions": ["PUT", "POST", "DELETE"]
 		},
 		"team": {
 		  "permissions": ["GET", "POST", "PUT"]
 		},
 		"v1-teams": {
 		  "permissions": ["DELETE"]
 		}
 	  }`
	// create a new reader with that JSON
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := PermissionSet{
		"team":     PermissionGroup{Permissions: []string{"GET", "POST", "PUT"}},
		"v1-teams": PermissionGroup{Permissions: []string{"DELETE"}},
		"v1-users": PermissionGroup{Permissions: []string{"PUT", "POST", "DELETE"}},
	}

	myPermissions, err := client.MyPermissions(getContext(nil))

	assert.Equal(t, expectedResponse, myPermissions)
	assert.Equal(t, nil, err)
}

func TestReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	myPermissions, err := client.MyPermissions(getContext(nil))

	assert.Equal(t, PermissionSet(PermissionSet(nil)), myPermissions)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/permissions",
		Method: http.MethodGet,
	}, err)
}

func TestReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	myPermissions, err := client.MyPermissions(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, PermissionSet(PermissionSet(nil)), myPermissions)
}

func TestHasPermissionStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyPermissions(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/permissions",
		Method: http.MethodGet,
	}, err)
}

func TestPermissionSetChecksPermission(t *testing.T) {
	permissions := PermissionSet{
		"user": {
			Permissions: []string{"GET", "PATCH"},
		},
		"team": {
			Permissions: []string{"GET"},
		},
	}

	assert.True(t, permissions.HasPermission("user", "PATCH"))
	assert.True(t, permissions.HasPermission("team", "GET"))
	assert.True(t, permissions.HasPermission("team", "get"))
	assert.False(t, permissions.HasPermission("team", "PATCHs"))
}
