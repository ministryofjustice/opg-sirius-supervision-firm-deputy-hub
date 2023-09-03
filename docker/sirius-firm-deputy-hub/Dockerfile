FROM golang:1.20.4-alpine3.16@sha256:6469405d7297f82d56195c90a3270b0806ef4bd897aa0628477d9959ab97a577 as base
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates

FROM node:19.9.0-alpine3.16@sha256:d73a7d0f2ec5d9c0f4b8aeb7768151e0a3283edba46f16adeec883ff703230cf as asset-env

WORKDIR /app

COPY web/assets web/assets
COPY package.json .
COPY yarn.lock .

RUN yarn install
RUN yarn build

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/cosmtrek/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-supervision-firm-deputy-hub

FROM alpine:3

WORKDIR /go/bin

RUN apk --update --no-cache add ca-certificates tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache libcrypto3 libssl3

COPY --from=build-env /go/bin/opg-sirius-supervision-firm-deputy-hub opg-sirius-supervision-firm-deputy-hub
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-supervision-firm-deputy-hub"]
