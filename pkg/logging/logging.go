package logging

import (
	"log/slog"
	"os"
)

func New(service string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(Service(service))
}
func Err(err error) slog.Attr { return slog.String("error", ErrorMessage(err)) }
