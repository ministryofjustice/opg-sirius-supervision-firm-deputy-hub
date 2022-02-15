package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type FirmHubDeputyTabInformation interface {
	GetFirmDeputies(sirius.Context, int) ([]sirius.FirmDeputy, error)
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type listDeputiesVars struct {
	Path                string
	XSRFToken           string
	FirmDeputiesDetails []sirius.FirmDeputy
	FirmDetails         sirius.FirmDetails
	Error               string
	ErrorMessage        string
	Errors              sirius.ValidationErrors
}

func renderTemplateForDeputyTab(client FirmHubDeputyTabInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])
		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		firmDeputiesDetails, err := client.GetFirmDeputies(ctx, firmId)
		if err != nil {
			return err
		}

		vars := listDeputiesVars{
			Path:                r.URL.Path,
			XSRFToken:           ctx.XSRFToken,
			FirmDeputiesDetails: firmDeputiesDetails,
			FirmDetails:         firmDetails,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
