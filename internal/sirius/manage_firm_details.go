package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) ManageFirmDetails(ctx Context, amendedFirmDetails model.FirmDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(model.FirmDetails{
		ID:           amendedFirmDetails.ID,
		FirmName:     amendedFirmDetails.FirmName,
		Email:        amendedFirmDetails.Email,
		PhoneNumber:  amendedFirmDetails.PhoneNumber,
		AddressLine1: amendedFirmDetails.AddressLine1,
		AddressLine2: amendedFirmDetails.AddressLine2,
		AddressLine3: amendedFirmDetails.AddressLine3,
		Town:         amendedFirmDetails.Town,
		County:       amendedFirmDetails.County,
		Postcode:     amendedFirmDetails.Postcode,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/firms/%d", amendedFirmDetails.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}
		return newStatusError(resp)
	}
	return err
}
