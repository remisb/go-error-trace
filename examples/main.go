package main

import (
	"log/slog"
	"os"

	trace "github.com/remisb/go-error-trace"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	// Usage examples
	err1 := trace.New("database connection failed")
	logger.Error("service initialization failed", slog.Any("error", err1))

	err2 := trace.Wrap(err1, "failed to start server")
	logger.Error("startup failed", slog.Any("error", err2))
}
