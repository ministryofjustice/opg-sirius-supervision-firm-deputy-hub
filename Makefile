all: go-lint unit-test build scan cypress down

.PHONY: cypress

build:
	docker compose build firm-deputy-hub

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build firm-deputy-hub yarn

clean:
	docker compose down
	docker compose run --rm yarn

compile-assets:
	docker compose run --rm yarn build

cypress: setup-directories clean
	docker compose up -d --wait firm-deputy-hub
	docker compose run --rm cypress run --env grepUntagged=true

dev-up:
	docker compose run --rm yarn
	docker compose run --rm yarn build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up --build firm-deputy-hub json-server

down:
	docker compose down

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-firm-deputy-hub:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-firm-deputy-hub:latest

setup-directories: test-results

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

unit-test: setup-directories
	docker compose run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

up: clean
	docker compose -f docker-compose.yml up --build firm-deputy-hub json-server yarn
