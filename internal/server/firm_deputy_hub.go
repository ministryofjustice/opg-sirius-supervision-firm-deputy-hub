package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type FirmDeputyHubInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmDeputyHubVars struct {
	Path      string
	XSRFToken string
	Error     string
	Errors    sirius.ValidationErrors
	FirmDetails sirius.FirmDetails
}

func renderTemplateForFirmHub(client FirmDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		url := r.URL.Path
		idFromParams := strings.Trim(url, "/")
		if idFromParams == "" {
			idFromParams = "0"
		}

		firmId, _ := strconv.Atoi(idFromParams)
		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil{
			return err
		}

		vars := firmDeputyHubVars{
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

