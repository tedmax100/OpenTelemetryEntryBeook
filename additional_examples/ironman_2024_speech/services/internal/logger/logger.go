package logger

import (
	"log/slog"
	"os"
	"time"

	slogformatter "github.com/samber/slog-formatter"
)

func GetLogger() *slog.Logger {
	logger := slog.New(
		slogformatter.NewFormatterHandler(
			slogformatter.TimezoneConverter(time.UTC),
			slogformatter.TimeFormatter(time.RFC3339, nil),
		)(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	)
	return logger
}

func (l *slog.Logger) With(key, value string) *slog.Logger {
	return l.With(key, value)
}
