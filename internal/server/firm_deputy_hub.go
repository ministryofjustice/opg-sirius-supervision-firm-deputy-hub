package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"net/http"
)

type FirmDeputyHubInformation interface {
}

type firmDeputyHubVars struct {
	Path      string
	XSRFToken string
	Error     string
	Errors    sirius.ValidationErrors
}

func renderTemplateForFirmDeputyHub(client FirmDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		vars := firmDeputyHubVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}