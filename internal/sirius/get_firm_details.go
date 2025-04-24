package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
	"time"
)

type ExecutiveCaseManagerOutgoing struct {
	EcmId int `json:"ecmId"`
}

func (c *Client) GetFirmDetails(ctx Context, firmId int) (model.FirmDetails, error) {
	var v model.FirmDetails

	requestURL := fmt.Sprintf(SupervisionAPIPath + "/v1/firms/%d", firmId)
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
	if v.PiiAmount != 0 {
		v.PiiAmountCommaFormat = humanize.Commaf(v.PiiAmount)
		v.PiiAmountIntFormat = int64(v.PiiAmount)
	}
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
