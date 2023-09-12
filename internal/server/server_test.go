package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"io"
)

type mockTemplates struct {
	count    int
	lastName string
	lastVars interface{}
	error    error
}

func (m *mockTemplates) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars
	return nil
}

type mockApiClient struct {
	error              error
	CurrentUserDetails model.Assignee
	FirmDetails        model.FirmDetails
}

func (m mockApiClient) GetUserDetails(sirius.Context) (model.Assignee, error) {
	return m.CurrentUserDetails, m.error
}

func (m mockApiClient) GetFirmDetails(sirius.Context, int) (model.FirmDetails, error) {
	return m.FirmDetails, m.error
}
