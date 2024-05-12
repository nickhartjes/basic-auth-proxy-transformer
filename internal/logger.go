package internal

import (
	"log/slog"
	"os"
)

// SetLogger sets the logger for the application
func SetLogger(settings Settings) {
	var level slog.Level
	if settings.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})
	logger := slog.New(logHandler)
	logger.Info(level.String() + " logging enabled")
	slog.SetDefault(logger)
}
