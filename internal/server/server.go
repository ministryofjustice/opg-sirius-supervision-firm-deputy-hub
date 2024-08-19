package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

//
//type Client interface {
//	ErrorHandlerClient
//	ManagePiiDetailsInformation
//	ManageFirmDetailsInformation
//	RequestPiiDetailsInformation
//	FirmHubDeputyTabInformation
//	ChangeECMClient
//}

type ApiClient interface {
	GetUserDetails(sirius.Context) (model.Assignee, error)
	GetFirmDetails(sirius.Context, int) (model.FirmDetails, error)
	EditPiiCertificate(sirius.Context, model.PiiDetails) error
	ManageFirmDetails(sirius.Context, model.FirmDetails) error
	RequestPiiCertificate(sirius.Context, sirius.PiiDetailsRequest) error
	GetFirmDeputies(sirius.Context, int) ([]model.FirmDeputy, error)
	GetProTeamUsers(sirius.Context) ([]model.TeamMembers, []model.Member, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, model.FirmDetails) error
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client ApiClient, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	//router := mux.NewRouter().StrictSlash(true)
	//router.Handle("/health-check", healthCheck())

	mux := http.NewServeMux()
	mux.Handle("GET /{firmId}/", wrap(renderTemplateForFirmHub(client, templates["firm-hub.gotmpl"])))

	mux.Handle("POST /{firmId}/", wrap(renderTemplateForFirmHub(client, templates["firm-hub.gotmpl"])))

	mux.Handle("GET /{firmId}/manage-pii-details", wrap(renderTemplateForManagePiiDetails(client, templates["manage-pii-details.gotmpl"])))
	mux.Handle("GET /{firmId}/manage-firm-details", wrap(renderTemplateForManageFirmDetails(client, templates["manage-firm-details.gotmpl"])))
	mux.Handle("GET /{firmId}/deputies", wrap(renderTemplateForDeputyTab(client, templates["deputies.gotmpl"])))
	mux.Handle("GET /{firmId}/health-check", healthCheck())
	mux.Handle("GET /{firmId}/request-pii-details", wrap(renderTemplateForRequestPiiDetails(client, templates["request-pii-details.gotmpl"])))
	mux.Handle("GET /{firmId}/change-ecm", wrap(renderTemplateForChangeECM(client, templates["change-ecm.gotmpl"])))

	//static := http.FileServer(http.Dir(envVars.WebDir))
	//mux.Handle("/assets/", static)
	//mux.Handle("/javascript/", static)
	//mux.Handle("/stylesheets/", static)

	static := staticFileHandler(envVars.WebDir)
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	//router.NotFoundHandler = wrap(notFoundHandler(templates["error.gotmpl"], envVars))

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, securityheaders.Use(mux)), "supervision-firm-hub")
}

func staticFileHandler(webDir string) http.Handler {
	h := http.FileServer(http.Dir(webDir + "/static"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
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
