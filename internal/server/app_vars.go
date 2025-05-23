package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AppVars struct {
	Path        string
	XSRFToken   string
	User        model.Assignee
	FirmDetails model.FirmDetails
	Error       string
	Errors      sirius.ValidationErrors
	EnvironmentVars
}

func (a AppVars) FirmId() int {
	return a.FirmDetails.ID
}

type AppVarsClient interface {
	GetUserDetails(sirius.Context) (model.Assignee, error)
	GetFirmDetails(sirius.Context, int) (model.FirmDetails, error)
}

func NewAppVars(client AppVarsClient, r *http.Request, envVars EnvironmentVars) (*AppVars, error) {
	ctx := getContext(r)
	group, groupCtx := errgroup.WithContext(ctx.Context)

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       ctx.XSRFToken,
		EnvironmentVars: envVars,
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
		firmId, _ := strconv.Atoi(r.PathValue("id"))
		firm, err := client.GetFirmDetails(ctx.With(groupCtx), firmId)
		if err != nil {
			return err
		}
		vars.FirmDetails = firm
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &vars, nil
}
