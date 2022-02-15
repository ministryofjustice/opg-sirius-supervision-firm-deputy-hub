package sirius

import (
	"encoding/json"
	"net/http"
)

type Member struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type TeamMembers struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Members     []Member `json:"members"`
}

func (c *Client) GetProTeamUsers(ctx Context) ([]TeamMembers, []Member, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams?type=pro", nil)
	if err != nil {
		return []TeamMembers{}, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return []TeamMembers{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []TeamMembers{}, nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return []TeamMembers{}, nil, newStatusError(resp)
	}

	var v []TeamMembers
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return []TeamMembers{}, nil, err
	}

	var ListOfAllProTeamMembers []Member

	for _, m := range v {
		for _, n := range m.Members {
			ListOfAllProTeamMembers = append(ListOfAllProTeamMembers, Member{
				Id:          n.Id,
				DisplayName: n.DisplayName,
			})
		}
	}

	return nil, ListOfAllProTeamMembers, nil
}
