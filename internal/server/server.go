package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

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
	ChangeECMInformation
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *logging.Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := wrapHandler(logger, client, templates["error.gotmpl"], prefix, siriusPublicURL)

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

	static := staticFileHandler(webDir)
	router.PathPrefix("/assets/").Handler(static)
	router.PathPrefix("/javascript/").Handler(static)
	router.PathPrefix("/stylesheets/").Handler(static)

	router.NotFoundHandler = notFoundHandler(templates["error.gotmpl"], siriusPublicURL)

	return http.StripPrefix(prefix, securityheaders.Use(router))
}

type Redirect string

func (e Redirect) Error() string {
	return "redirect to " + string(e)
}

func (e Redirect) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type AppVars struct {
	Path      string
	XSRFToken string
	User      sirius.Assignee
	Firm      sirius.FirmDetails
	Error     string
	Errors    sirius.ValidationErrors
}

type Handler func(app AppVars, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Code      int
	Error     string
	Errors    bool
}

type ErrorHandlerClient interface {
	GetUserDetails(sirius.Context) (sirius.Assignee, error)
	GetFirmDetails(sirius.Context, int) (sirius.FirmDetails, error)
}

func wrapHandler(logger *logging.Logger, client ErrorHandlerClient, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := getContext(r)

			group, groupCtx := errgroup.WithContext(ctx.Context)

			vars := AppVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
			}

			group.Go(func() error {
				user, err := client.GetUserDetails(ctx.With(groupCtx))
				if err != nil {
					return err
				}
				vars.User = user
				return nil
			})
			group.Go(func() error {
				firmId, _ := strconv.Atoi(mux.Vars(r)["id"])
				firm, err := client.GetFirmDetails(ctx.With(groupCtx), firmId)
				if err != nil {
					return err
				}
				vars.Firm = firm
				return nil
			})

			err := group.Wait()

			if err == nil {
				err = next(vars, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}

func notFoundHandler(tmplError Template, siriusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmplError.ExecuteTemplate(w, "page", errorVars{
			SiriusURL: siriusURL,
			Code:      http.StatusNotFound,
			Error:     "Not Found",
		})
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
