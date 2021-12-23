package sirius

import (
	"encoding/json"
	"net/http"
	"time"
)

type Deputy struct {
	DeputyId         int    `json:"id"`
	DeputyNumber     int    `json:"deputyNumber"`
	OrganisationName string `json:"organisationName"`
}

type FirmDetails struct {
	ID                    int      `json:"id"`
	FirmName              string   `json:"firmName"`
	FirmNumber            int      `json:"firmNumber"`
	Email                 string   `json:"email"`
	PhoneNumber           string   `json:"phoneNumber"`
	AddressLine1          string   `json:"addressLine1"`
	AddressLine2          string   `json:"addressLine2"`
	AddressLine3          string   `json:"addressLine3"`
	Town                  string   `json:"town"`
	County                string   `json:"county"`
	Postcode              string   `json:"postcode"`
	Deputies              []Deputy `json:"deputies"`
	PiiExpiry             string   `json:"piiExpiry"`
	PiiAmount             string   `json:"piiAmount"`
	TotalNumberOfDeputies int
}

func (c *Client) GetFirmDetails(ctx Context, firmId int) (FirmDetails, error) {
	var v FirmDetails

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/firms/1", nil)

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

	v.PiiExpiry = reformatDate(v.PiiExpiry)

	return v, err
}

func reformatDate(dateString string) string {
	if dateString == "" {
		return ""
	}
	dateTime, _ := time.Parse("2006-01-02", dateString)

	return dateTime.Format("02/01/2006")
}
