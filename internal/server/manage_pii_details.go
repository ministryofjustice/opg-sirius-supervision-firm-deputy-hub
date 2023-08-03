package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type ManagePiiDetailsInformation interface {
	EditPiiCertificate(sirius.Context, sirius.PiiDetails) error
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmHubManagePiiVars struct {
	FirmDetails          sirius.FirmDetails
	ErrorMessage         string
	AddFirmPiiDetailForm sirius.PiiDetails
	AppVars
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		vars := firmHubManagePiiVars{
			FirmDetails: firmDetails,
		}
		vars.AppVars = app

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			addFirmPiiDetailForm := sirius.PiiDetails{
				FirmId:      firmId,
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

			err = client.EditPiiCertificate(ctx, addFirmPiiDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.AddFirmPiiDetailForm = addFirmPiiDetailForm
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d?success=piiDetails", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
