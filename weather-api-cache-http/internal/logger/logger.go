package logger

import (
	"os"
	"log/slog"
)

func New() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
}
