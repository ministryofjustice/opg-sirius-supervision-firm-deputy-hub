package model

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
