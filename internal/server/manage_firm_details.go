package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type ManageFirmDetailsInformation interface {
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

type firmHubManageFirmVars struct {
	Path                 string
	XSRFToken            string
	Error                string
	Errors               sirius.ValidationErrors
	FirmDetails          sirius.FirmDetails
	ErrorMessage         string
	AddFirmPiiDetailForm sirius.PiiDetails
}

func renderTemplateForManageFirmDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
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

		//case http.MethodPost:
		//	addFirmPiiDetailForm := sirius.PiiDetails{
		//		FirmId:       firmId,
		//		PiiReceived:  r.PostFormValue("pii-received"),
		//		PiiExpiry:    r.PostFormValue("pii-expiry"),
		//		PiiRequested: r.PostFormValue("pii-requested"),
		//	}
		//
		//	if r.PostFormValue("pii-amount") != "" {
		//		piiAmountFloat, err := strconv.ParseFloat(r.PostFormValue("pii-amount"), 64)
		//		if err != nil {
		//			return err
		//		}
		//		addFirmPiiDetailForm.PiiAmount = piiAmountFloat
		//	}
		//
		//	err = client.EditPiiCertificate(ctx, addFirmPiiDetailForm)
		//
		//	if verr, ok := err.(sirius.ValidationError); ok {
		//		verr.Errors = renameEditPiiValidationErrorMessages(verr.Errors)
		//		vars := firmHubManagePiiVars{
		//			Path:                 r.URL.Path,
		//			XSRFToken:            ctx.XSRFToken,
		//			Errors:               verr.Errors,
		//			FirmDetails:          firmDetails,
		//			AddFirmPiiDetailForm: addFirmPiiDetailForm,
		//		}
		//		return tmpl.ExecuteTemplate(w, "page", vars)
		//	}
		//
		//	return Redirect(fmt.Sprintf("/%d?success=piiDetails", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
