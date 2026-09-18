package model

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
type FirmResponse struct {
	ID                   int                  `json:"id"`
	FirmName             string               `json:"firmName"`
	FirmNumber           int                  `json:"firmNumber"`
	Email                string               `json:"email"`
	PhoneNumber          string               `json:"phoneNumber"`
	AddressLine1         string               `json:"addressLine1"`
	AddressLine2         string               `json:"addressLine2"`
	AddressLine3         string               `json:"addressLine3"`
	Town                 string               `json:"town"`
	County               string               `json:"county"`
	Postcode             string               `json:"postcode"`
	ExecutiveCaseManager ExecutiveCaseManager `json:"executiveCaseManager"`
	Deputies             []DeputyResponse     `json:"deputies"`
	PiiReceived          string               `json:"piiReceived"`
	PiiExpiry            string               `json:"piiExpiry"`
	PiiAmount            float64              `json:"piiAmount,omitempty"`
	PiiRequested         string               `json:"piiRequested"`
}

type FirmDetails struct {
	ID                     int                  `json:"id"`
	FirmName               string               `json:"firmName"`
	FirmNumber             int                  `json:"firmNumber"`
	Email                  string               `json:"email"`
	PhoneNumber            string               `json:"phoneNumber"`
	AddressLine1           string               `json:"addressLine1"`
	AddressLine2           string               `json:"addressLine2"`
	AddressLine3           string               `json:"addressLine3"`
	Town                   string               `json:"town"`
	County                 string               `json:"county"`
	Postcode               string               `json:"postcode"`
	ExecutiveCaseManager   ExecutiveCaseManager `json:"executiveCaseManager"`
	Deputies               []FirmDeputies       `json:"deputies"`
	PiiReceived            string               `json:"piiReceived"`
	PiiExpiry              string               `json:"piiExpiry"`
	PiiAmount              float64              `json:"piiAmount,omitempty"`
	PiiRequested           string               `json:"piiRequested"`
	PiiReceivedDateFormat  string
	PiiExpiryDateFormat    string
	PiiRequestedDateFormat string
	TotalNumberOfDeputies  int
	PiiAmountCommaFormat   string
	PiiAmountIntFormat     int64
}

type FirmDeputies struct {
	DeputyId         int    `json:"id"`
	DeputyNumber     int    `json:"deputyNumber"`
	OrganisationName string `json:"organisationName"`
}
