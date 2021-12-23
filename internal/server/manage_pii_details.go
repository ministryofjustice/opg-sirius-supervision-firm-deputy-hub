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
	Path                 string
	XSRFToken            string
	Error                string
	Errors               sirius.ValidationErrors
	FirmDetails          sirius.FirmDetails
	ErrorMessage         string
	AddFirmPiiDetailForm sirius.PiiDetails
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
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

			vars := firmHubManagePiiVars{
				Path:        r.URL.Path,
				XSRFToken:   ctx.XSRFToken,
				FirmDetails: firmDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			addFirmPiiDetailForm := sirius.PiiDetails{
				FirmId:       firmId,
				PiiReceived:  r.PostFormValue("pii-received"),
				PiiExpiry:    r.PostFormValue("pii-expiry"),
				PiiAmount:    r.PostFormValue("pii-amount"),
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			fmt.Println(r.PostFormValue("pii-received"))

			err := client.EditPiiCertificate(ctx, addFirmPiiDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditPiiValidationErrorMessages(verr.Errors)
				vars := firmHubManagePiiVars{
					Path:                 r.URL.Path,
					XSRFToken:            ctx.XSRFToken,
					Errors:               verr.Errors,
					FirmDetails:          firmDetails,
					AddFirmPiiDetailForm: addFirmPiiDetailForm,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/%d?success=piiDetails", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameEditPiiValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if fieldName == "piiReceived" && errorType == "isEmpty" {
				err[errorType] = "The PII received date is required and can't be empty"
				errorCollection["pii-received"] = err
			} else if fieldName == "piiExpiry" && errorType == "isEmpty" {
				err[errorType] = "The PII expiry date is required and can't be empty"
				errorCollection["pii-expiry"] = err
			} else if fieldName == "piiAmount" && errorType == "isEmpty" {
				err[errorType] = "The PII amount is required and can't be empty"
				errorCollection["pii-amount"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
