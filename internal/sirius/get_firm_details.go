package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return FirmDetails{}, ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return FirmDetails{ID: 0, FirmName: "", FirmNumber: 0}, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
