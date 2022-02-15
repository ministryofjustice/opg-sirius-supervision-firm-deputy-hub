package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Deputy struct {
	DeputyId         int    `json:"id"`
	DeputyNumber     int    `json:"deputyNumber"`
	OrganisationName string `json:"organisationName"`
}

type ExecutiveCaseManager struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

type ExecutiveCaseManagerOutgoing struct {
	EcmId int `json:"ecmId"`
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
	Deputies               []Deputy             `json:"deputies"`
	PiiReceived            string               `json:"piiReceived"`
	PiiExpiry              string               `json:"piiExpiry"`
	PiiAmount              float64              `json:"piiAmount,omitempty"`
	PiiRequested           string               `json:"piiRequested"`
	PiiReceivedDateFormat  string
	PiiExpiryDateFormat    string
	PiiRequestedDateFormat string
	TotalNumberOfDeputies  int
	PiiAmountCommaFormat   string
}

func (c *Client) GetFirmDetails(ctx Context, firmId int) (FirmDetails, error) {
	var v FirmDetails

	requestURL := fmt.Sprintf("/api/v1/firms/%d", firmId)
	req, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)

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

	v.PiiReceivedDateFormat = reformatDatesForAutofill(v.PiiReceived)
	v.PiiExpiryDateFormat = reformatDatesForAutofill(v.PiiExpiry)
	v.PiiRequestedDateFormat = reformatDatesForAutofill(v.PiiRequested)

	return v, err
}

func reformatDatesForAutofill(date string) string {
	if date != "" {
		dateTime, _ := time.Parse("02/01/2006", date)
		date = dateTime.Format("2006-01-02")
		return date
	}
	return ""
}
