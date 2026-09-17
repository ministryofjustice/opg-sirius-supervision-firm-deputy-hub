package model

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
