package server

import (
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
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

func New(logger *logging.Logger, client Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/health-check", healthCheck())

	pageRouter := router.PathPrefix("/{id}").Subrouter()
	pageRouter.Use(logging.Use(logger))

	pageRouter.Handle("",
		wrap(
			renderTemplateForFirmHub(client, templates["firm-hub.gotmpl"])))

	pageRouter.Handle("/manage-pii-details",
		wrap(
			renderTemplateForManagePiiDetails(client, templates["manage-pii-details.gotmpl"])))

	pageRouter.Handle("/manage-firm-details",
		wrap(
			renderTemplateForManageFirmDetails(client, templates["manage-firm-details.gotmpl"])))

	pageRouter.Handle("/deputies",
		wrap(
			renderTemplateForDeputyTab(client, templates["deputies.gotmpl"])))

	router.Handle("/health-check", healthCheck())

	pageRouter.Handle("/request-pii-details",
		wrap(
			renderTemplateForRequestPiiDetails(client, templates["request-pii-details.gotmpl"])))

	pageRouter.Handle("/change-ecm",
		wrap(
			renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))

	static := staticFileHandler(envVars.WebDir)
	router.PathPrefix("/assets/").Handler(static)
	router.PathPrefix("/javascript/").Handler(static)
	router.PathPrefix("/stylesheets/").Handler(static)

	router.NotFoundHandler = wrap(notFoundHandler(templates["error.gotmpl"], envVars))

	return http.StripPrefix(envVars.Prefix, securityheaders.Use(router))
}

func notFoundHandler(tmplError Template, envVars EnvironmentVars) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		_ = tmplError.ExecuteTemplate(w, "page", ErrorVars{
			Code:            http.StatusNotFound,
			Error:           "Page not found",
			EnvironmentVars: envVars,
		})
		return nil
	}
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
	h := http.FileServer(http.Dir(webDir + "/static"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
}
