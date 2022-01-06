package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRequestPiiDetailsInformation struct {
	count       int
	lastCtx     sirius.Context
	err         error
	firmDetails sirius.FirmDetails
}

func (m *mockRequestPiiDetailsInformation) RequestPiiCertificate(ctx sirius.Context, piiData sirius.PiiDetailsRequest) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockRequestPiiDetailsInformation) GetFirmDetails(ctx sirius.Context, firmId int) (sirius.FirmDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmDetails, m.err
}

func TestRequestPiiDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestPiiDetailsInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForRequestPiiDetails(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostRequestPii(t *testing.T) {
	assert := assert.New(t)
	client := &mockRequestPiiDetailsInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForRequestPiiDetails(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123?success=requestPiiDetails"))
}

func TestErrorRequestPiiMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockRequestPiiDetailsInformation{}

	validationErrors := sirius.ValidationErrors{
		"piiRequested": {
			"isEmpty": "The PII requested date is required and can't be empty",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForRequestPiiDetails(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"piiRequested": {
			"isEmpty": "The PII requested date is required and can't be empty",
			},
		},
	}

	assert.Equal(expectedValidationErrors, returnedError)
}
