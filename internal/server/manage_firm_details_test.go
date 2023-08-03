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

type mockManageFirmDetailsInformation struct {
	count       int
	lastCtx     sirius.Context
	err         error
	firmDetails sirius.FirmDetails
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
	err := handler(AppVars{}, w, r)

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
		returnedError = renderTemplateForManageFirmDetails(client, nil)(AppVars{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/123?success=firmDetails"))
}
