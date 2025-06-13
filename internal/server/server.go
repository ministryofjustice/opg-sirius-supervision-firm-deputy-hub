package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

type Client interface {
	ErrorHandlerClient
	FirmHubInformation
	ManagePiiDetailsInformation
	ManageFirmDetailsInformation
	RequestPiiDetailsInformation
	FirmHubDeputyTabInformation
	ChangeECMClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	mux := http.NewServeMux()
	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	// Static file routes
	static := staticFileHandler(envVars.WebDir)
	mux.Handle("/static/assets/", static)
	mux.Handle("/static/javascript/", static)
	mux.Handle("/static/stylesheets/", static)

	// Main routes
	mux.Handle("/{id}", wrap(renderTemplateForFirmHub(client, templates["firm-hub.gotmpl"])))
	mux.Handle("/{id}/manage-pii-details", wrap(renderTemplateForManagePiiDetails(client, templates["manage-pii-details.gotmpl"])))
	mux.Handle("/{id}/manage-firm-details", wrap(renderTemplateForManageFirmDetails(client, templates["manage-firm-details.gotmpl"])))
	mux.Handle("/{id}/deputies", wrap(renderTemplateForDeputyTab(client, templates["deputies.gotmpl"])))
	mux.Handle("/{id}/request-pii-details", wrap(renderTemplateForRequestPiiDetails(client, templates["request-pii-details.gotmpl"])))
	mux.Handle("/{id}/change-ecm", wrap(renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))

	// Health check
	mux.Handle("/health-check", healthCheck())

	// 404 fallback
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = templates["error.gotmpl"].ExecuteTemplate(w, "page", ErrorVars{
			Code:            http.StatusNotFound,
			Error:           "Page not found",
			EnvironmentVars: envVars,
		})
	})

	return http.StripPrefix(envVars.Prefix, securityheaders.Use(telemetry.Middleware(logger)(mux)))
}

func getContext(r *http.Request) sirius.Context {
	token := ""
	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}
	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}

func staticFileHandler(webDir string) http.Handler {
	h := http.FileServer(http.Dir(webDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
}
