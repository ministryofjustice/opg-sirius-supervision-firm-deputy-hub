package model

import "strings"

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

type RAGRating struct {
	Name   string
	Colour string
}

func (fd FirmDeputy) GetRAGRating() RAGRating {
	var rag RAGRating
	switch strings.ToUpper(fd.MarkedAsClass) {
	case "RED":
		rag.Name = "High risk"
		rag.Colour = "red"
	case "AMBER":
		rag.Name = "Medium risk"
		rag.Colour = "orange"
	case "GREEN":
		rag.Name = "Low risk"
		rag.Colour = "green"
	}
	return rag
}
