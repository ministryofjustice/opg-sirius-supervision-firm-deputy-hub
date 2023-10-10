package model

type DeputyImportantInformation struct {
	PanelDeputy bool `json:"panelDeputy"`
}

type Deputies struct {
	DeputyId                   int                        `json:"id"`
	Firstname                  string                     `json:"firstname"`
	Surname                    string                     `json:"surname"`
	DeputyNumber               int                        `json:"deputyNumber"`
	Orders                     []Orders                   `json:"orders"`
	ExecutiveCaseManager       ExecutiveCaseManager       `json:"executiveCaseManager"`
	OrganisationName           string                     `json:"organisationName"`
	Town                       string                     `json:"town"`
	Assurance                  Assurance                  `json:"mostRecentlyCompletedAssurance"`
	DeputyImportantInformation DeputyImportantInformation `json:"deputyImportantInformation"`
}

type FirmDeputy struct {
	DeputyId             int
	Firstname            string
	Surname              string
	DeputyNumber         int
	ActiveClientsCount   int
	ExecutiveCaseManager string
	OrganisationName     string
	Town                 string
	ReviewDate           string
	MarkedAsLabel        string
	MarkedAsClass        string
	AssuranceType        string
	PanelDeputy          bool
}
