package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ChangeECM(ctx Context, changeECMForm ExecutiveCaseManagerOutgoing, firmDetails FirmDetails) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(ExecutiveCaseManagerOutgoing{EcmId: changeECMForm.EcmId})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/firms/%d/ecm", firmDetails.ID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
