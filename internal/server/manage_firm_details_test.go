package server

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestManageFirmDetails(t *testing.T) {
	assert := assert.New(t)

	client := mockClient
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := renderTemplateForManageFirmDetails(client, template)(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostManageFirm(t *testing.T) {
	client := mockClient

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	returnedError := renderTemplateForManageFirmDetails(client, nil)(AppVars{FirmDetails: mockClient.firmDetails}, w, r)

	assert.Equal(t, Redirect("/firm/123?success=firmDetails"), returnedError)
}

func TestPostManageFirmReturnsError(t *testing.T) {
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
			client := mockClient
			client.manageFirmDetailsErr = test.apiError
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			err := renderTemplateForManageFirmDetails(client, template)(AppVars{}, w, r)

			if test.wantValidationErrors != nil {
				assert.Equal(t, test.apiError.(sirius.ValidationError).Errors, template.lastVars.(firmHubManageFirmVars).Errors)
				assert.Equal(t, test.wantCode, w.Result().StatusCode)
			} else {
				assert.Nil(t, template.lastVars)
				assert.Equal(t, test.wantError, err)
			}
		})
	}
}
