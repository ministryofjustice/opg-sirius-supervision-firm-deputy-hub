all: go-lint unit-test build scan acceptance-testing cypress

.PHONY: cypress

go-lint:
	docker compose -f docker/docker-compose.ci.yml run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache

setup-directories: test-results

unit-test: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

build:
	docker compose -f docker/docker-compose.ci.yml build firm-deputy-hub

scan:
	trivy image sirius/sirius-firm-deputy-hub:latest

cypress:
	docker compose -f docker/docker-compose.ci.yml run cypress
