package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	PiiAmount             float64  `json:"piiAmount",omitempty`
	TotalNumberOfDeputies int
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

	return v, err
}
