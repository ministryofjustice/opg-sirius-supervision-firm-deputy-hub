package util

import (
	"log/slog"
	"net/http"
	"os"
)

func NewLogger(serviceName string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "level" {
					return slog.Attr{}
				}

				if a.Key == "time" {
					a.Key = "timestamp"
				}

				if a.Key == "msg" {
					a.Key = "message"
				}

				return a
			},
		}).WithAttrs([]slog.Attr{slog.String("service_name", serviceName)}))
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
