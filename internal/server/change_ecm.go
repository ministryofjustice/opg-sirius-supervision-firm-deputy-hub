package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type ChangeECMInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
	GetProTeamUsers(sirius.Context) ([]sirius.TeamMembers, []sirius.Member, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.FirmDetails) error
}

type changeECMHubVars struct {
	Path           string
	XSRFToken      string
	FirmDetails    sirius.FirmDetails
	EcmTeamDetails []sirius.Member
	Error          string
	Errors         sirius.ValidationErrors
	AppVars
}

func renderTemplateForChangeECM(client ChangeECMInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		_, ecmTeamDetails, err := client.GetProTeamUsers(ctx)
		if err != nil {
			return err
		}

		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		vars := changeECMHubVars{
			Path:           r.URL.Path,
			XSRFToken:      ctx.XSRFToken,
			FirmDetails:    firmDetails,
			EcmTeamDetails: ecmTeamDetails,
		}
		vars.AppVars = app

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			EcmIdStringValue := r.PostFormValue("select-ecm")

			if EcmIdStringValue == "" {
				vars.Errors = sirius.ValidationErrors{
					"Change ECM": {"": "Select an executive case manager"},
				}
				EcmIdStringValue = "0"
			}

			EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			err = client.ChangeECM(ctx, changeECMForm, firmDetails)

			if len(vars.Errors) >= 1 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			return Redirect(fmt.Sprintf("/%d?success=ecm", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
