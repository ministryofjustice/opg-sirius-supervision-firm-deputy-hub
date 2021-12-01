package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Deputies struct {
	DeputyId int `json:"id"`
	DeputyNumber int `json:"deputyNumber"`
	OrganisationName string `json:"organisationName"`
}

type FirmDetails struct {
	ID                               int                  `json:"id"`
	FirmName                 		string             	  `json:"firmName"`
	FirmNumber                 		int               	`json:"firmNumber"`
	Email                            string               `json:"email"`
	PhoneNumber                      string               `json:"phoneNumber"`
	AddressLine1                     string               `json:"addressLine1"`
	AddressLine2                     string               `json:"addressLine2"`
	AddressLine3                     string               `json:"addressLine3"`
	Town                             string               `json:"town"`
	County                           string               `json:"county"`
	Postcode                         string               `json:"postcode"`
	Deputies 						[]Deputies 			`json:"deputies"`
	TotalNumberOfDeputies int
}


func (c *Client) GetFirmDetails(ctx Context, firmId int) (FirmDetails, error) {
	var v FirmDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/firm/%d", firmId), nil)

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

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var ve struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return v, ValidationError{
				Errors: ve.ValidationErrors,
			}
		}

		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	v.TotalNumberOfDeputies = calculateNumberOfDeputies(v.Deputies)

	return v, err
}

func calculateNumberOfDeputies(deputyArray []Deputies) int {
	totalDeputies := 0
	for range deputyArray {
		totalDeputies++
	}
	return totalDeputies
}