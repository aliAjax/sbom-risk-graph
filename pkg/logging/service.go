package logging

import (
	"log/slog"
	"strings"
)

func Service(name string) slog.Attr {
	return slog.String("service", strings.TrimSpace(name))
}
