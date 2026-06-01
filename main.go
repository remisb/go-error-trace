package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// StackTraceError is a custom error that captures a stack trace.
type StackTraceError struct {
	err   error
	stack []uintptr
}

// Error implements the error interface.
func (e *StackTraceError) Error() string {
	return e.err.Error()
}

// Unwrap supports errors.Is and errors.As.
func (e *StackTraceError) Unwrap() error {
	return e.err
}

// StackTrace returns the captured stack for formatting.
func (e *StackTraceError) StackTrace() []uintptr {
	return e.stack
}

// New creates a new error with stack trace.
func New(msg string) error {
	return &StackTraceError{
		err:   errors.New(msg),
		stack: captureStack(3), // Skip New + runtime frames
	}
}

// Wrap wraps an existing error and adds a stack trace.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &StackTraceError{
		err:   fmt.Errorf("%s: %w", msg, err),
		stack: captureStack(3),
	}
}

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceErrorWithStack,
	})

	logger := slog.New(handler)

	// Usage examples
	err1 := New("database connection failed")
	logger.Error("service initialization failed", slog.Any("error", err1))

	err2 := Wrap(err1, "failed to start server")
	logger.Error("startup failed", slog.Any("error", err2))
}

// replaceErrorWithStack expands errors that have stack traces.
func replaceErrorWithStack(_ []string, a slog.Attr) slog.Attr {
	if err, ok := a.Value.Any().(error); ok {
		var ste *StackTraceError
		if errors.As(err, &ste) {
			return slog.Attr{
				Key: a.Key,
				Value: slog.GroupValue(
					slog.String("msg", ste.Error()),
					slog.Any("trace", formatStack(ste.StackTrace())),
				),
			}
		}
		// Fallback for normal errors
		return slog.String(a.Key, err.Error())
	}
	return a
}

func formatStack(pcs []uintptr) []string {
	frames := runtime.CallersFrames(pcs)
	var trace []string

	for {
		frame, more := frames.Next()
		// Skip runtime and slog internal frames if desired
		if !strings.Contains(frame.File, "runtime/") {
			trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}
	return trace
}

func captureStack(skip int) []uintptr {
	const depth = 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs)
	return pcs[:n]
}
