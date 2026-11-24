package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"log/slog"
	"net/http"
	"net/url"
)

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

type Handler func(app AppVars, w http.ResponseWriter, r *http.Request) error

type ErrorVars struct {
	Code  int
	Error string
	EnvironmentVars
}

type ExpandedError interface {
	Title() string
	Data() interface{}
}

func LoggerRequest(l *slog.Logger, r *http.Request, err error) {
	if ee, ok := err.(ExpandedError); ok {
		l.Info("request error",
			slog.String("title", ee.Title()),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()),
			slog.Any("data", ee.Data()))
	} else if err != nil {
		l.Info("request error",
			slog.String("error", err.Error()),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	} else {
		l.Info("",
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	}
}

type ErrorHandlerClient interface {
	GetUserDetails(sirius.Context) (model.Assignee, error)
	GetFirmDetails(sirius.Context, int) (model.FirmDetails, error)
}

func wrapHandler(logger *slog.Logger, client ErrorHandlerClient, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars, err := NewAppVars(client, r, envVars)

			if err == nil {
				err = next(*vars, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					redirect := ""

					if r.RequestURI != "" {
						redirect = "?redirect=" + url.QueryEscape(r.RequestURI)
					}

					http.Redirect(w, r, envVars.SiriusPublicURL+"/auth"+redirect, http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					u, err := url.Parse(redirect.To())
					if err != nil || u.Host != "" || u.Scheme != "" {
						http.Error(w, "Invalid redirect", http.StatusBadRequest)
						return
					}
					// Use only the path component to prevent open redirect
					safeRedirect := u.RequestURI()
					http.Redirect(w, r, envVars.Prefix+safeRedirect, http.StatusFound) // #nosec G710 -- host and scheme rejected above
					return
				}

				LoggerRequest(logger, r, err)

				code := http.StatusInternalServerError
				if serverStatusError, ok := err.(StatusError); ok {
					code = serverStatusError.Code()
				}
				if siriusStatusError, ok := err.(sirius.StatusError); ok {
					code = siriusStatusError.Code
				}

				w.WriteHeader(code)
				errVars := ErrorVars{
					Code:            code,
					Error:           err.Error(),
					EnvironmentVars: envVars,
				}
				err = tmplError.ExecuteTemplate(w, "page", errVars)

				if err != nil {
					LoggerRequest(logger, r, nil)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
