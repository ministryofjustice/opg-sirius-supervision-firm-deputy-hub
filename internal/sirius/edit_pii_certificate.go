package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
)

func (c *Client) EditPiiCertificate(ctx Context, editPiiData model.PiiDetails) error {
	var k model.PiiDetails
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(model.PiiDetails{
		PiiReceived:  editPiiData.PiiReceived,
		PiiExpiry:    editPiiData.PiiExpiry,
		PiiAmount:    editPiiData.PiiAmount,
		PiiRequested: editPiiData.PiiRequested,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/firms/%d/indemnity-insurance", editPiiData.FirmId)

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

	err = json.NewDecoder(resp.Body).Decode(&k)
	return err
}
