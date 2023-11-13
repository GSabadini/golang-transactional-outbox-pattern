package logger

import (
	"log/slog"
	"os"
)

var Slog = slog.New(slog.NewJSONHandler(os.Stdout, nil))
