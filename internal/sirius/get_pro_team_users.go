package sirius

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) GetProTeamUsers(ctx Context) ([]model.TeamMembers, []model.Member, error) {
	req, err := c.newRequest(ctx, http.MethodGet, SupervisionAPIPath+"/v1/teams?type=pro", nil)
	if err != nil {
		return []model.TeamMembers{}, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return []model.TeamMembers{}, nil, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return []model.TeamMembers{}, nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return []model.TeamMembers{}, nil, newStatusError(resp)
	}

	var v []model.TeamMembers
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return []model.TeamMembers{}, nil, err
	}

	var ListOfAllProTeamMembers []model.Member

	for _, m := range v {
		for _, n := range m.Members {
			ListOfAllProTeamMembers = append(ListOfAllProTeamMembers, model.Member{
				Id:          n.Id,
				DisplayName: n.DisplayName,
			})
		}
	}

	return nil, ListOfAllProTeamMembers, nil
}
