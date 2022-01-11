package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManageFirmDetailsInformation struct {
	count       int
	lastCtx     sirius.Context
	err         error
	firmDetails sirius.FirmDetails
}

func (m *mockManageFirmDetailsInformation) GetFirmDetails(ctx sirius.Context, firmId int) (sirius.FirmDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmDetails, m.err
}

func (m *mockManageFirmDetailsInformation) ManageFirmDetails(ctx sirius.Context, firmDetails sirius.FirmDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestManageFirmDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageFirmDetailsInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForManageFirmDetails(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostManageFirm(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageFirmDetailsInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManageFirmDetails(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123?success=firmDetails"))
}

//func TestErrorManagePiiMessageWhenIsEmpty(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockManagePiiDetailsInformation{}
//
//	validationErrors := sirius.ValidationErrors{
//		"piiReceived": {
//			"isEmpty": "The PII received date is required and can't be empty",
//		},
//		"piiExpiry": {
//			"isEmpty": "The PII expiry is required and can't be empty",
//		},
//		"piiAmount": {
//			"isEmpty": "The PII amount is required and can't be empty",
//		},
//	}
//
//	client.err = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
//	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//
//	var returnedError error
//
//	testHandler := mux.NewRouter()
//	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
//		returnedError = renderTemplateForManagePiiDetails(client, template)(sirius.PermissionSet{}, w, r)
//	})
//
//	testHandler.ServeHTTP(w, r)
//
//	expectedValidationErrors := sirius.ValidationError{
//		Errors: sirius.ValidationErrors{
//			"piiReceived": {
//				"isEmpty": "The PII received date is required and can't be empty",
//			},
//			"piiExpiry": {
//				"isEmpty": "The PII expiry is required and can't be empty",
//			},
//			"piiAmount": {
//				"isEmpty": "The PII amount is required and can't be empty",
//			},
//		},
//	}
//
//	assert.Equal(expectedValidationErrors, returnedError)
//}
