package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
)

type RequestPiiDetailsInformation interface {
	RequestPiiCertificate(sirius.Context, sirius.PiiDetailsRequest) error
}

type firmHubRequestPiiVars struct {
	ErrorMessage          string
	RequestPiiDetailsForm sirius.PiiDetailsRequest
	AppVars
}

func renderTemplateForRequestPiiDetails(client RequestPiiDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		vars := firmHubRequestPiiVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			requestPiiDetailsForm := sirius.PiiDetailsRequest{
				FirmId:       app.FirmId(),
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			err := client.RequestPiiCertificate(ctx, requestPiiDetailsForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				vars.RequestPiiDetailsForm = requestPiiDetailsForm
				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=requestPiiDetails", app.FirmId()))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
