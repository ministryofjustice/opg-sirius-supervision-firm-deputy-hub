package main

import (
	"context"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/util"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/server"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
)

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-firm-deputy-hub ")

	envVars, err := server.NewEnvironmentVars()
	if err != nil {
		logger.Fatal(err)
	}

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"prefix": func(s string) string {
				return envVars.Prefix + s
			},
			"sirius": func(s string) string {
				return envVars.SiriusPublicURL + s
			},
			"prohub": func(s string) string {
				return envVars.ProHubURL + s
			},
			"rename_errors": util.RenameErrors,
		}).
		ParseGlob(envVars.WebDir + "/template/*/*.gotmpl")

	files, _ := filepath.Glob(envVars.WebDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	client, err := sirius.NewClient(http.DefaultClient, envVars.SiriusURL)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":" + envVars.Port,
		Handler: server.New(logger, client, tmpls, envVars),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + envVars.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}
