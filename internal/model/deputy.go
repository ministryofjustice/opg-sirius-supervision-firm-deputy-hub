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

type DeputyResponse struct {
	AddressLine1          string `json:"addressLine1"`
	AddressLine2          string `json:"addressLine2"`
	AddressLine3          string `json:"addressLine3"`
	CorrespondenceByEmail bool   `json:"correspondenceByEmail"`
	CorrespondenceByPhone bool   `json:"correspondenceByPhone"`
	CorrespondenceByPost  bool   `json:"correspondenceByPost"`
	CorrespondenceByWelsh bool   `json:"correspondenceByWelsh"`
	Country               string `json:"country"`
	County                string `json:"county"`
	DateOfDeath           string `json:"dateOfDeath"`
	DeputyNumber          int    `json:"deputyNumber"`
	DeputyStatus          string `json:"deputyStatus"`
	DeputySubType         struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"deputySubType"`
	DeputyType struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"deputyType"`
	Dob                  string `json:"dob"`
	Email                string `json:"email"`
	EveningNumber        string `json:"eveningNumber"`
	ExecutiveCaseManager struct {
		DisplayName string `json:"displayName"`
		Id          int    `json:"id"`
		Name        string `json:"name"`
	} `json:"executiveCaseManager"`
	Id                                int             `json:"id"`
	IsAirmailRequired                 bool            `json:"isAirmailRequired"`
	MobileNumber                      string          `json:"mobileNumber"`
	NormalizedUid                     int64           `json:"normalizedUid"`
	Orders                            [][]interface{} `json:"orders"`
	OrganisationName                  string          `json:"organisationName"`
	OrganisationTeamOrDepartmentName  string          `json:"organisationTeamOrDepartmentName"`
	PersonType                        string          `json:"personType"`
	Postcode                          string          `json:"postcode"`
	SpecialCorrespondenceRequirements struct {
		AudioTape                  bool `json:"audioTape"`
		HearingImpaired            bool `json:"hearingImpaired"`
		LargePrint                 bool `json:"largePrint"`
		SpellingOfNameRequiresCare bool `json:"spellingOfNameRequiresCare"`
	} `json:"specialCorrespondenceRequirements"`
	Town string `json:"town"`
	UId  string `json:"uId"`
}
