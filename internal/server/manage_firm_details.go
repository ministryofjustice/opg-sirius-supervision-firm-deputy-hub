package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type ManageFirmDetailsInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
	ManageFirmDetails(sirius.Context, sirius.FirmDetails) error
}

type firmHubManageFirmVars struct {
	Path                string
	XSRFToken           string
	Error               string
	Errors              sirius.ValidationErrors
	FirmDetails         sirius.FirmDetails
	ErrorMessage        string
	EditFirmDetailsForm sirius.FirmDetails
	AppVars
}

func renderTemplateForManageFirmDetails(client ManageFirmDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		vars := firmHubManageFirmVars{
			Path:        r.URL.Path,
			XSRFToken:   ctx.XSRFToken,
			FirmDetails: firmDetails,
		}
		vars.AppVars = app

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			editFirmDetailsForm := sirius.FirmDetails{
				ID:           firmId,
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

			err = client.ManageFirmDetails(ctx, editFirmDetailsForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.EditFirmDetailsForm = editFirmDetailsForm
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d?success=firmDetails", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
