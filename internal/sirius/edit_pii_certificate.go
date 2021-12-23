package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PiiDetails struct {
	FirmId       int    `json:"firmId"`
	PiiReceived  string `json:"piiReceived"`
	PiiExpiry    string `json:"piiExpiry"`
	PiiAmount    string `json:"piiAmount"`
	PiiRequested string `json:"piiRequested"`
}

func (c *Client) EditPiiCertificate(ctx Context, editPiiData PiiDetails) error {
	var k PiiDetails
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(PiiDetails{
		PiiReceived:  editPiiData.PiiReceived,
		PiiExpiry:    editPiiData.PiiExpiry,
		PiiAmount:    editPiiData.PiiAmount,
		PiiRequested: editPiiData.PiiRequested,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/firms/%d/pii", editPiiData.FirmId)

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

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
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

	err = json.NewDecoder(resp.Body).Decode(&k)
	return err
}
