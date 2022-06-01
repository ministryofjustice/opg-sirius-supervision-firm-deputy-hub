package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type RequestPiiDetailsInformation interface {
	RequestPiiCertificate(sirius.Context, sirius.PiiDetailsRequest) error
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmHubRequestPiiVars struct {
	Path                  string
	XSRFToken             string
	Error                 string
	Errors                sirius.ValidationErrors
	FirmDetails           sirius.FirmDetails
	ErrorMessage          string
	RequestPiiDetailsForm sirius.PiiDetailsRequest
}

func renderTemplateForRequestPiiDetails(client RequestPiiDetailsInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		firmDetails, err := client.GetFirmDetails(ctx, firmId)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			vars := firmHubRequestPiiVars{
				Path:        r.URL.Path,
				XSRFToken:   ctx.XSRFToken,
				FirmDetails: firmDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			requestPiiDetailsForm := sirius.PiiDetailsRequest{
				FirmId:       firmId,
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			err = client.RequestPiiCertificate(ctx, requestPiiDetailsForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars := firmHubRequestPiiVars{
					Path:                  r.URL.Path,
					XSRFToken:             ctx.XSRFToken,
					Errors:                verr.Errors,
					FirmDetails:           firmDetails,
					RequestPiiDetailsForm: requestPiiDetailsForm,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d?success=requestPiiDetails", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
