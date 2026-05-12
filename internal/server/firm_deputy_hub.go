package server

import (
	"html"
	"html/template"
	"net/http"
	"strings"
)

type FirmHubInformation interface {
}

type firmHubVars struct {
	SuccessMessage template.HTML
	AppVars
}

func renderTemplateForFirmHub(client FirmHubInformation, tmpl Template) Handler {
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

		switch splitString[1] {
		case "firm":
			return "Firm changed to " + html.EscapeString(firmName)
		case "newFirm":
			return "Firm added"
		case "firmDetails":
			return "Firm details updated"
		case "deputyDetails":
			return "Deputy details updated"
		case "piiDetails":
			return "PII details updated"
		case "requestPiiDetails":
			return "PII details requested"
		case "ecm":
			return "<abbr title='Executive Case Manager'>ECM</abbr> changed to " + html.EscapeString(ecmName)
		}
	}
	return ""
}
