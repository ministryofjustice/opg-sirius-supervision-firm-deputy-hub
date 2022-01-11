package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type FirmHubInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmHubVars struct {
	Path           string
	XSRFToken      string
	Error          string
	Errors         sirius.ValidationErrors
	FirmDetails    sirius.FirmDetails
	SuccessMessage string
}

func renderTemplateForFirmHub(client FirmHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		url := r.URL.Path
		idFromParams := strings.Trim(url, "/")

		firmId, _ := strconv.Atoi(idFromParams)
		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		successMessage := createSuccessAndSuccessMessageForVars(r.URL.String(), firmDetails.FirmName)

		vars := firmHubVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			FirmDetails:    firmDetails,
			SuccessMessage: successMessage,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func createSuccessAndSuccessMessageForVars(url, firmName string) string {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")

		if splitString[1] == "firm" {
			return "Firm changed to " + firmName
		} else if splitString[1] == "newFirm" {
			return "Firm added"
		} else if splitString[1] == "firmDetails" {
				return "Firm details updated"
		} else if splitString[1] == "deputyDetails" {
			return "Deputy details updated"
		} else if splitString[1] == "piiDetails" {
			return "PII details updated"
		}

	}
	return ""
}
