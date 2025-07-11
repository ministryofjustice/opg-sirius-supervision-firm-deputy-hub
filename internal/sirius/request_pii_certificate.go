package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PiiDetailsRequest struct {
	FirmId       int    `json:"firmId"`
	PiiRequested string `json:"piiRequested"`
}

func (c *Client) RequestPiiCertificate(ctx Context, requestPiiData PiiDetailsRequest) error {
	var k PiiDetailsRequest
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(PiiDetailsRequest{
		FirmId:       requestPiiData.FirmId,
		PiiRequested: requestPiiData.PiiRequested,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/firms/%d/indemnity-insurance", requestPiiData.FirmId)

	req, err := c.newRequest(ctx, http.MethodPatch, requestURL, &body)

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

	err = json.NewDecoder(resp.Body).Decode(&k)
	return err
}
