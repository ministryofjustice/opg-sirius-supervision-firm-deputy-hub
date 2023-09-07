package server

import (
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockChangeECMClient struct {
	count                int
	lastCtx              sirius.Context
	getProTeamUsersError error
	changeECMError       error
	EcmTeamDetails       []model.Member
	EcmTeamApiDetails    []model.TeamMembers
}

func (m *mockChangeECMClient) GetProTeamUsers(ctx sirius.Context) ([]model.TeamMembers, []model.Member, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.EcmTeamApiDetails, m.EcmTeamDetails, m.getProTeamUsersError
}

func (m *mockChangeECMClient) ChangeECM(ctx sirius.Context, changeEcmForm sirius.ExecutiveCaseManagerOutgoing, firmDetails model.FirmDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.changeECMError
}

func TestGetChangeECM(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangeECMClient{}
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
	client := &mockChangeECMClient{}

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

func TestPostChangeECMReturnsError(t *testing.T) {
	tests := []struct {
		apiError             error
		wantValidationErrors sirius.ValidationErrors
		wantError            error
		wantCode             int
	}{
		{
			apiError: sirius.ValidationError{
				Errors: sirius.ValidationErrors{
					"Error": {"": "Test error"},
				},
			},
			wantValidationErrors: sirius.ValidationErrors{
				"Error": {"": "Test error"},
			},
			wantError: nil,
			wantCode:  400,
		},
		{
			apiError:             sirius.StatusError{Code: 503},
			wantValidationErrors: nil,
			wantError:            sirius.StatusError{Code: 503},
			wantCode:             503,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			client := &mockChangeECMClient{changeECMError: test.apiError}
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/76/firm-ecm", strings.NewReader("select-ecm=1"))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			err := renderTemplateForChangeECM(client, template)(AppVars{}, w, r)

			if test.wantValidationErrors != nil {
				assert.Equal(t, test.apiError.(sirius.ValidationError).Errors, template.lastVars.(changeECMHubVars).Errors)
				assert.Equal(t, test.wantCode, w.Result().StatusCode)
			} else {
				assert.Nil(t, template.lastVars)
				assert.Equal(t, test.wantError, err)
			}
		})
	}
}
