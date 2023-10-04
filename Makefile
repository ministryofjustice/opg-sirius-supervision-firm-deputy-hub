all: go-lint unit-test build scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

unit-test: setup-directories
	docker compose run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

build:
	docker compose build firm-deputy-hub

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-firm-deputy-hub:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-firm-deputy-hub:latest

cypress: setup-directories
	docker compose up -d --wait firm-deputy-hub
	docker compose run --rm cypress run --env grepUntagged=true

up:
	docker compose up --build -d firm-deputy-hub

dev-up:
	docker compose run --rm yarn
	docker compose run --rm yarn build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up --build firm-deputy-hub json-server

down:
	docker compose down
