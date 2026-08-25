package logging

import (
	"log/slog"
	"os"
)

func New(service string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", service)
}
func Err(err error) slog.Attr { return slog.String("error", err.Error()) }
