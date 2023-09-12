package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
)

type ChangeECMClient interface {
	GetProTeamUsers(sirius.Context) ([]model.TeamMembers, []model.Member, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, model.FirmDetails) error
}

type changeECMHubVars struct {
	EcmTeamDetails []model.Member
	AppVars
}

func renderTemplateForChangeECM(client ChangeECMClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		_, ecmTeamDetails, err := client.GetProTeamUsers(ctx)
		if err != nil {
			return err
		}

		vars := changeECMHubVars{
			EcmTeamDetails: ecmTeamDetails,
			AppVars:        app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			EcmIdStringValue := r.PostFormValue("select-ecm")

			if EcmIdStringValue == "" {
				vars.Errors = sirius.ValidationErrors{
					"Change ECM": {"": "Select an executive case manager"},
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			err = client.ChangeECM(ctx, changeECMForm, app.FirmDetails)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=ecm", app.FirmId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
