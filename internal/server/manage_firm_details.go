package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
)

type ManageFirmDetailsInformation interface {
	ManageFirmDetails(sirius.Context, sirius.FirmDetails) error
}

type firmHubManageFirmVars struct {
	ErrorMessage        string
	EditFirmDetailsForm sirius.FirmDetails
	AppVars
}

func renderTemplateForManageFirmDetails(client ManageFirmDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)

		vars := firmHubManageFirmVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			editFirmDetailsForm := sirius.FirmDetails{
				ID:           app.FirmId(),
				FirmName:     r.PostFormValue("firm-name"),
				Email:        r.PostFormValue("email"),
				PhoneNumber:  r.PostFormValue("telephone"),
				AddressLine1: r.PostFormValue("address-line-1"),
				AddressLine2: r.PostFormValue("address-line-2"),
				AddressLine3: r.PostFormValue("address-line-3"),
				Town:         r.PostFormValue("town"),
				County:       r.PostFormValue("county"),
				Postcode:     r.PostFormValue("postcode"),
			}

			err := client.ManageFirmDetails(ctx, editFirmDetailsForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.EditFirmDetailsForm = editFirmDetailsForm
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=firmDetails", app.FirmId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
