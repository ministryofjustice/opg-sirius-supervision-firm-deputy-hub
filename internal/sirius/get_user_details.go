package sirius

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Assignee struct {
	ID       int      `json:"id"`
	Roles    []string `json:"roles"`
	Username string   `json:"displayName"`
}

func (a Assignee) GetRoles() string {
	return strings.Join(a.Roles, ",")
}

func (c *Client) GetUserDetails(ctx Context) (Assignee, error) {
	var v Assignee

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/users/current", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return v, err
}
