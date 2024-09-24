package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type ManagePiiDetailsInformation interface {
	EditPiiCertificate(sirius.Context, model.PiiDetails) error
}

type firmHubManagePiiVars struct {
	ErrorMessage         string
	AddFirmPiiDetailForm model.PiiDetails
	AppVars
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := firmHubManagePiiVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			addFirmPiiDetailForm := model.PiiDetails{
				FirmId:      app.FirmId(),
				PiiReceived: r.PostFormValue("pii-received"),
				PiiExpiry:   r.PostFormValue("pii-expiry"),
			}

			if r.PostFormValue("pii-amount") != "" {
				piiAmountFloat, err := strconv.ParseFloat(r.PostFormValue("pii-amount"), 64)
				if err != nil {
					return err
				}
				addFirmPiiDetailForm.PiiAmount = piiAmountFloat
			}

			err := client.EditPiiCertificate(ctx, addFirmPiiDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.AddFirmPiiDetailForm = addFirmPiiDetailForm
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/firm/%d?success=piiDetails", app.FirmId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
