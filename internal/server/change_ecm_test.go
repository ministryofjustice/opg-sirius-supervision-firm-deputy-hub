package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMInformation struct {
	count             int
	lastCtx           sirius.Context
	err               error
	EcmTeamDetails    []model.Member
	EcmTeamApiDetails []model.TeamMembers
}

func (m *mockChangeECMInformation) GetProTeamUsers(ctx sirius.Context) ([]model.TeamMembers, []model.Member, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.EcmTeamApiDetails, m.EcmTeamDetails, m.err
}

func (m *mockChangeECMInformation) ChangeECM(ctx sirius.Context, changeEcmForm sirius.ExecutiveCaseManagerOutgoing, firmDetails model.FirmDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestGetChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForChangeECM(client, template)
	app := AppVars{FirmDetails: mockFirmDetails}
	err := handler(app, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changeECMHubVars{
		AppVars: app,
	}, template.lastVars)
}

func TestPostChangeECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/firm-ecm", strings.NewReader("{ecmId:26}"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/firm-ecm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeECM(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(returnedError)
}

func TestPostChangeECMReturnsErrorWithNoECM(t *testing.T) {
	assert := assert.New(t)
	client := &mockChangeECMInformation{}

	validationErrors := sirius.ValidationErrors{
		"Change ECM": {"": "Select an executive case manager"},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/76/firm-ecm", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}/firm-ecm", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForChangeECM(client, template)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationError := sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"Change ECM": {
				"": "Select an executive case manager",
			},
		},
	}

	assert.Equal(expectedValidationError, returnedError)
}
