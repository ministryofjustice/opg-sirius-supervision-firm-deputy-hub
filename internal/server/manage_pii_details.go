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
}

type firmHubManagePiiVars struct {
	Path         string
	XSRFToken    string
	Error        string
	Errors       sirius.ValidationErrors
	ErrorMessage string
	FirmId       int
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		fmt.Println("firm id")
		fmt.Println(firmId)

		switch r.Method {
		case http.MethodGet:

			vars := firmHubManagePiiVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
				FirmId:    firmId,
			}

			fmt.Println("get method firm hub manage vars")
			fmt.Println(vars)

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:

			addFirmPiiDetailForm := sirius.PiiDetails{
				FirmId:       firmId,
				PiiReceived:  r.PostFormValue("pii-received"),
				PiiExpiry:    r.PostFormValue("pii-expiry"),
				PiiAmount:    r.PostFormValue("pii-amount"),
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			err := client.EditPiiCertificate(ctx, addFirmPiiDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditPiiValidationErrorMessages(verr.Errors)
				vars := firmHubManagePiiVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return RedirectError(fmt.Sprintf("/%d?success=piiDetails", firmId))

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
				err[errorType] = "The pii received date is required and can't be empty"
				errorCollection["piiReceived"] = err
			} else if fieldName == "piiExpiry" && errorType == "isEmpty" {
				err[errorType] = "The pii expiry is required and can't be empty"
				errorCollection["piiExpiry"] = err
			} else if fieldName == "piiAmount" && errorType == "isEmpty" {
				err[errorType] = "The pii amount is required and can't be empty"
				errorCollection["piiAmount"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
