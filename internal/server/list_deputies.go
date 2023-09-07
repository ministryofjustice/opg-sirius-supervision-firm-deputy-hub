package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
)

type FirmHubDeputyTabInformation interface {
	GetFirmDeputies(sirius.Context, int) ([]model.FirmDeputy, error)
}

type listDeputiesVars struct {
	FirmDeputiesDetails []model.FirmDeputy
	ErrorMessage        string
	AppVars
}

func renderTemplateForDeputyTab(client FirmHubDeputyTabInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		firmDeputiesDetails, err := client.GetFirmDeputies(ctx, app.FirmId())
		if err != nil {
			return err
		}

		vars := listDeputiesVars{
			FirmDeputiesDetails: firmDeputiesDetails,
			AppVars:             app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
