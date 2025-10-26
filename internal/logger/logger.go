package logger

import (
	"log/slog"
	"os"
)

// InitLogger initializes and returns a configured slog.Logger
func InitLogger() *slog.Logger {
	// Create a text handler that writes to stdout
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	return slog.New(handler)
}
