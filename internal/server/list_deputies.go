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
	FirmDeputiesDetails []sirius.FirmDeputy
	FirmDetails         sirius.FirmDetails
	ErrorMessage        string
	AppVars
}

func renderTemplateForDeputyTab(client FirmHubDeputyTabInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
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
			FirmDeputiesDetails: firmDeputiesDetails,
			FirmDetails:         firmDetails,
		}
		vars.AppVars = app

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
