package model

import "strings"

type Assignee struct {
	ID       int      `json:"id"`
	Roles    []string `json:"roles"`
	Username string   `json:"displayName"`
}

func (a Assignee) GetRoles() string {
	return strings.Join(a.Roles, ",")
}
