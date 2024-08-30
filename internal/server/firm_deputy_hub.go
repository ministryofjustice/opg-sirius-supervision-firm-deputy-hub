package server

import (
	"html/template"
	"net/http"
	"strings"
)

type firmHubVars struct {
	SuccessMessage template.HTML
	AppVars
}

func renderTemplateForFirmHub(client ApiClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		successMessage := createSuccessAndSuccessMessageForVars(r.URL.String(), app.FirmDetails.FirmName, app.FirmDetails.ExecutiveCaseManager.DisplayName)

		vars := firmHubVars{
			SuccessMessage: template.HTML(successMessage),
			AppVars:        app,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func createSuccessAndSuccessMessageForVars(url, firmName, ecmName string) string {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")

		if splitString[1] == "firm" {
			return "Firm changed to " + firmName
		} else if splitString[1] == "newFirm" {
			return "Firm added"
		} else if splitString[1] == "firmDetails" {
			return "Firm details updated"
		} else if splitString[1] == "deputyDetails" {
			return "Deputy details updated"
		} else if splitString[1] == "piiDetails" {
			return "PII details updated"
		} else if splitString[1] == "requestPiiDetails" {
			return "PII details requested"
		} else if splitString[1] == "ecm" {
			return "<abbr title='Executive Case Manager'>ECM</abbr> changed to " + ecmName
		}

	}
	return ""
}
