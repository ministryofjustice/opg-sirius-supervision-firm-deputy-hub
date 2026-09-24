package model

type OrderStatus struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type Order struct {
	Id          int         `json:"id"`
	Client      Client      `json:"client"`
	OrderStatus OrderStatus `json:"orderStatus"`
}

type Orders struct {
	Order Order `json:"order"`
}

// "client": {
// "addressLine1": matchers.Like(""),
// "addressLine2": matchers.Like(""),
// "addressLine3": matchers.Like(""),
// "correspondenceByEmail": matchers.Like(false),
// "correspondenceByPhone": matchers.Like(false),
// "correspondenceByPost": matchers.Like(false),
// "correspondenceByWelsh": matchers.Like(false),
// "country": matchers.Like(""),
// "county": matchers.Like(""),
// "dateOfDeath": matchers.Like(""),
// "dob": matchers.Like(""),
// "email": matchers.Like(""),
// "firstname": matchers.Like("tKI7Z94H2VQHfAM"),
// "id": matchers.Like(88),
// "interpreterRequired": matchers.Like(""),
// "isAirmailRequired": matchers.Like(false),
// "middlenames": matchers.Like(""),
// "normalizedUid": matchers.Like(700000003256),
// "postcode": matchers.Like(""),
// "salutation": matchers.Like(""),
// "specialCorrespondenceRequirements": {
// "audioTape": matchers.Like(false),
// "hearingImpaired": matchers.Like(false),
// "largePrint": matchers.Like(false),
// "spellingOfNameRequiresCare": matchers.Like(false),
// },
// "surname": matchers.Like("SRuHE3Hj7APfZd1"),
// "town": matchers.Like(""),
// "uId": matchers.Like("7000-0000-3256"),
// },
type T struct {
	Order struct {
		Client struct {
			AddressLine1                      string `json:"addressLine1"`
			AddressLine2                      string `json:"addressLine2"`
			AddressLine3                      string `json:"addressLine3"`
			CorrespondenceByEmail             bool   `json:"correspondenceByEmail"`
			CorrespondenceByPhone             bool   `json:"correspondenceByPhone"`
			CorrespondenceByPost              bool   `json:"correspondenceByPost"`
			CorrespondenceByWelsh             bool   `json:"correspondenceByWelsh"`
			Country                           string `json:"country"`
			County                            string `json:"county"`
			DateOfDeath                       string `json:"dateOfDeath"`
			Dob                               string `json:"dob"`
			Email                             string `json:"email"`
			Firstname                         string `json:"firstname"`
			Id                                int    `json:"id"`
			InterpreterRequired               string `json:"interpreterRequired"`
			IsAirmailRequired                 bool   `json:"isAirmailRequired"`
			Middlenames                       string `json:"middlenames"`
			NormalizedUid                     int64  `json:"normalizedUid"`
			Postcode                          string `json:"postcode"`
			Salutation                        string `json:"salutation"`
			SpecialCorrespondenceRequirements struct {
				AudioTape                  bool `json:"audioTape"`
				HearingImpaired            bool `json:"hearingImpaired"`
				LargePrint                 bool `json:"largePrint"`
				SpellingOfNameRequiresCare bool `json:"spellingOfNameRequiresCare"`
			} `json:"specialCorrespondenceRequirements"`
			Surname string `json:"surname"`
			Town    string `json:"town"`
			UId     string `json:"uId"`
		} `json:"client"`
		Id          int `json:"id"`
		OrderStatus struct {
			Handle string `json:"handle"`
			Label  string `json:"label"`
		} `json:"orderStatus"`
	} `json:"order"`
}
