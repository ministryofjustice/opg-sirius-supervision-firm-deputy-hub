package model

//type T struct {
//	Id       int `json:"id"`
//	Deputies []struct {
//		Id                                int    `json:"id"`
//		UId                               string `json:"uId"`
//		NormalizedUid                     int64  `json:"normalizedUid"`
//		Email                             string `json:"email"`
//		Dob                               string `json:"dob"`
//		DateOfDeath                       string `json:"dateOfDeath"`
//		Salutation                        string `json:"salutation"`
//		Firstname                         string `json:"firstname"`
//		Middlenames                       string `json:"middlenames"`
//		Surname                           string `json:"surname"`
//		AddressLine1                      string `json:"addressLine1"`
//		AddressLine2                      string `json:"addressLine2"`
//		AddressLine3                      string `json:"addressLine3"`
//		Town                              string `json:"town"`
//		County                            string `json:"county"`
//		Postcode                          string `json:"postcode"`
//		Country                           string `json:"country"`
//		IsAirmailRequired                 bool   `json:"isAirmailRequired"`
//		CorrespondenceByPost              bool   `json:"correspondenceByPost"`
//		CorrespondenceByPhone             bool   `json:"correspondenceByPhone"`
//		CorrespondenceByEmail             bool   `json:"correspondenceByEmail"`
//		CorrespondenceByWelsh             bool   `json:"correspondenceByWelsh"`
//		PersonType                        string `json:"personType"`
//		SpecialCorrespondenceRequirements struct {
//			AudioTape                  bool `json:"audioTape"`
//			LargePrint                 bool `json:"largePrint"`
//			HearingImpaired            bool `json:"hearingImpaired"`
//			SpellingOfNameRequiresCare bool `json:"spellingOfNameRequiresCare"`
//		} `json:"specialCorrespondenceRequirements"`
//		DeputyStatus  string          `json:"deputyStatus"`
//		Orders        [][]interface{} `json:"orders"`
//		MobileNumber  string          `json:"mobileNumber"`
//		EveningNumber string          `json:"eveningNumber"`
//		DeputyType    struct {
//			Handle string `json:"handle"`
//			Label  string `json:"label"`
//		} `json:"deputyType"`
//		DeputyNumber                     int    `json:"deputyNumber"`
//		OrganisationName                 string `json:"organisationName"`
//		OrganisationTeamOrDepartmentName string `json:"organisationTeamOrDepartmentName"`
//		ExecutiveCaseManager             struct {
//			Id          int    `json:"id"`
//			Name        string `json:"name"`
//			DisplayName string `json:"displayName"`
//		} `json:"executiveCaseManager"`
//		DeputySubType struct {
//			Handle string `json:"handle"`
//			Label  string `json:"label"`
//		} `json:"deputySubType"`
//	} `json:"deputies"`
//	FirmName     string `json:"firmName"`
//	AddressLine1 string `json:"addressLine1"`
//	AddressLine2 string `json:"addressLine2"`
//	AddressLine3 string `json:"addressLine3"`
//	Town         string `json:"town"`
//	County       string `json:"county"`
//	Postcode     string `json:"postcode"`
//	PhoneNumber  string `json:"phoneNumber"`
//	Email        string `json:"email"`
//	FirmNumber   int    `json:"firmNumber"`
//	PersonType   string `json:"personType"`
//}

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
