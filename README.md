# OPG SIRIUS SUPERVISION FIRM DEPUTY HUB

### Major dependencies

-   [Go](https://golang.org/) (>= 1.22)
-   [docker compose](https://docs.docker.com/compose/install/) (>= 2.0.0)

#### Installing dependencies locally:
(This is only necessary if dunning without docker)
-   `yarn install`
-   `go mod download`

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/deputies/firm/1`.

To enable debugging and hot-reloading of Go files:

`docker compose up --build firm-deputy-hub` or `make up`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates)
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

### Without docker

Alternatively to set it up not using Docker use below. This may be necessary to build the assets folder locally (if
there are assets missing) as the developer version of the docker compose file does not pass the Air stage. This hosts it on `localhost:1234`

-   `yarn install && yarn build ` #run this to build your assets folder locally
-   `go build main.go `
-   `./main `

    ***

## Run Cypress tests

`make cypress`

---

## Run the unit/functional tests

`make unit-test`

-------------------------------------------------------------------
## Run Trivy scanning

`make scan`
