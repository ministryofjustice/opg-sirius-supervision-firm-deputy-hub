package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockAppVarsClient struct {
	lastCtx sirius.Context
	err     error
	user    model.Assignee
	firm    model.FirmDetails
}

func (m *mockAppVarsClient) GetUserDetails(ctx sirius.Context) (model.Assignee, error) {
	m.lastCtx = ctx

	return m.user, m.err
}

func (m *mockAppVarsClient) GetFirmDetails(ctx sirius.Context, firmId int) (model.FirmDetails, error) {
	m.lastCtx = ctx

	return m.firm, m.err
}

var mockUserDetails = model.Assignee{
	ID: 1,
}

var mockFirmDetails = model.FirmDetails{
	ID: 123,
}

func TestNewAppVars(t *testing.T) {
	client := &mockAppVarsClient{user: mockUserDetails, firm: mockFirmDetails}
	r, _ := http.NewRequest("GET", "/path", nil)

	envVars := EnvironmentVars{}
	vars, err := NewAppVars(client, r, envVars)

	assert.Nil(t, err)
	assert.Equal(t, AppVars{
		Path:            "/path",
		XSRFToken:       "",
		User:            mockUserDetails,
		FirmDetails:     mockFirmDetails,
		Error:           "",
		Errors:          nil,
		EnvironmentVars: envVars,
	}, *vars)
}
