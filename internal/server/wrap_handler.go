package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
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

type ErrorVars struct {
	Code  int
	Error string
	EnvironmentVars
}

type ErrorHandlerClient interface {
	GetUserDetails(sirius.Context) (model.Assignee, error)
	GetFirmDetails(sirius.Context, int) (model.FirmDetails, error)
}

func wrapHandler(logger *logging.Logger, client ErrorHandlerClient, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
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
						redirect = "?redirect=" + r.RequestURI
					}

					http.Redirect(w, r, envVars.SiriusPublicURL+"/auth"+redirect, http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, envVars.Prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

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
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
