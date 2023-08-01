package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockFirmHubDeputyTabInformation struct {
	count               int
	lastCtx             sirius.Context
	err                 error
	firmDeputiesDetails []sirius.FirmDeputy
	firmDetails         sirius.FirmDetails
}

func (m *mockFirmHubDeputyTabInformation) GetFirmDeputies(ctx sirius.Context, firmId int) ([]sirius.FirmDeputy, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmDeputiesDetails, m.err
}

func (m *mockFirmHubDeputyTabInformation) GetFirmDetails(ctx sirius.Context, firmId int) (sirius.FirmDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmDetails, m.err
}

func TestRenderTemplateForDeputyTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockFirmHubDeputyTabInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyTab(client, template)
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(2, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}
