package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
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
	firmDeputiesDetails []model.FirmDeputy
}

func (m *mockFirmHubDeputyTabInformation) GetFirmDeputies(ctx sirius.Context, firmId int) ([]model.FirmDeputy, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firmDeputiesDetails, m.err
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

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}
