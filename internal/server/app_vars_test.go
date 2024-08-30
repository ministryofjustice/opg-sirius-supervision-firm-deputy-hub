package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var mockClient = mockApiClient{
	currentUserDetails: model.Assignee{
		ID:       99,
		Roles:    []string{"System Admin", "COP User"},
		Username: "test user",
	},
	firmDetails: model.FirmDetails{
		ID:       123,
		FirmName: "Legit Firm Inc",
	},
}

func TestNewAppVars(t *testing.T) {
	client := mockClient
	r, _ := http.NewRequest("GET", "/path", nil)

	envVars := EnvironmentVars{}
	vars, err := NewAppVars(client, r, envVars)

	assert.Nil(t, err)
	assert.Equal(t, AppVars{
		Path:            "/path",
		XSRFToken:       "",
		User:            mockClient.currentUserDetails,
		FirmDetails:     mockClient.firmDetails,
		Error:           "",
		Errors:          nil,
		EnvironmentVars: envVars,
	}, vars)
}
