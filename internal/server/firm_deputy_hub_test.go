package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanRenderFirmDetailsPage(t *testing.T) {
	assert := assert.New(t)

	client := mockClient
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/supervision/deputies/firm/3", nil)

	handler := renderTemplateForFirmHub(client, template)
	app := AppVars{FirmDetails: mockClient.firmDetails}
	err := handler(app, w, r)

	assert.Nil(err)
	assert.Equal("page", template.lastName)
	assert.Equal(firmHubVars{
		AppVars: app,
	}, template.lastVars)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
