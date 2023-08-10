package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockFirmHubInformation struct {
	count       int
	lastCtx     sirius.Context
	err         error
	firmDetails sirius.FirmDetails
}

func TestCanRenderFirmDetailsPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockFirmHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/supervision/deputies/firm/3", nil)

	handler := renderTemplateForFirmHub(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)
	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal("page", template.lastName)
	assert.Equal(firmHubVars{
		AppVars: AppVars{Path: "/supervision/deputies/firm/3"},
	}, template.lastVars)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
