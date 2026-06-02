package main

import (
	"log/slog"
	"os"

	trace "github.com/remisb/go-error-trace"
)

func main() {
	//handler := newTextHandler()
	handler := newJSONHandler()
	logger := slog.New(handler)

	// Usage examples
	err := errInFunction()
	logger.Error("service initialization failed", slog.Any("error", err))

	err2 := trace.Wrap(err, "failed to start server")
	logger.Error("startup failed", slog.Any("error", err2))
}

func newTextHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
}

func newJSONHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
}

func errInFunction() error {
	err := trace.New("database connection failed")
	return err
}
