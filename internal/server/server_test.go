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
	error                    error
	currentUserDetails       model.Assignee
	firmDetails              model.FirmDetails
	editPiiCertificateErr    error
	manageFirmDetailsErr     error
	requestPiiCertificateErr error
	firmDeputies             []model.FirmDeputy
	firmDeputiesErr          error
	proTeamUsersTeamMembers  []model.TeamMembers
	proTeamUsersMembers      []model.Member
	proTeamUsersErr          error
	changeEcmErr             error
}

func (m mockApiClient) GetUserDetails(sirius.Context) (model.Assignee, error) {
	return m.currentUserDetails, m.error
}

func (m mockApiClient) GetFirmDetails(sirius.Context, int) (model.FirmDetails, error) {
	return m.firmDetails, m.error
}

func (m mockApiClient) EditPiiCertificate(sirius.Context, model.PiiDetails) error {
	return m.editPiiCertificateErr
}

func (m mockApiClient) ManageFirmDetails(sirius.Context, model.FirmDetails) error {
	return m.manageFirmDetailsErr
}

func (m mockApiClient) RequestPiiCertificate(sirius.Context, sirius.PiiDetailsRequest) error {
	return m.requestPiiCertificateErr
}

func (m mockApiClient) GetFirmDeputies(sirius.Context, int) ([]model.FirmDeputy, error) {
	return m.firmDeputies, m.firmDeputiesErr
}

func (m mockApiClient) GetProTeamUsers(sirius.Context) ([]model.TeamMembers, []model.Member, error) {
	return m.proTeamUsersTeamMembers, m.proTeamUsersMembers, m.proTeamUsersErr
}

func (m mockApiClient) ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, model.FirmDetails) error {
	return m.changeEcmErr
}
