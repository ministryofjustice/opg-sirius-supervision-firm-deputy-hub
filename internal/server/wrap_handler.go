package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"log/slog"
	"net/http"
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

//type Handler interface {
//	render(app AppVars, w http.ResponseWriter, r *http.Request) error
//}

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
		l.Info(ee.Title(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()),
			slog.Any("data", ee.Data()))
	} else if err != nil {
		l.Info(err.Error(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	} else {
		l.Info("",
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	}
}

//type ErrorHandlerClient interface {
//	GetUserDetails(sirius.Context) (model.Assignee, error)
//	GetFirmDetails(sirius.Context, int) (model.FirmDetails, error)
//}

func wrapHandler(logger *slog.Logger, client ApiClient, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var err error
			vars, err := NewAppVars(client, r, envVars)
			if err == nil {
				err = next(vars, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					redirect := ""

					if r.RequestURI != "" {
						redirect = "?redirect=" + r.RequestURI
					}

					http.Redirect(w, r, envVars.SiriusPublicURL+"/auth"+redirect, http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, envVars.Prefix+redirect.To(), http.StatusFound)
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
