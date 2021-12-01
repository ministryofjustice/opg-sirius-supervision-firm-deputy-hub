package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type FirmHubInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmHubVars struct {
	Path      string
	XSRFToken string
	Error     string
	Errors    sirius.ValidationErrors
	FirmDetails sirius.FirmDetails
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
		if err != nil{
			return err
		}

		vars := firmHubVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
			FirmDetails: firmDetails,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

