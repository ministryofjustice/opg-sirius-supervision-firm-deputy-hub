FROM golang:1.16 as build-env

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-supervision-firm-deputy-hub

FROM alpine:3.14

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata

COPY --from=build-env /go/bin/opg-sirius-supervision-firm-deputy-hub opg-sirius-supervision-firm-deputy-hub
ENTRYPOINT ["./opg-sirius-supervision-firm-deputy-hub"]