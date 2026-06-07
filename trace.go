package trace

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

const defaultStackDepth = 32
const filterOutPackageOn = true
const filterOutPackage = "runtime"

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

// LogValue implements slog.LogValuer for automatic structured logging.
func (e *StackTraceError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("msg", e.err.Error()),
		slog.Any("trace", formatStack(e.stack)),
	)
}

// New creates a new error with stack trace.
func New(msg string) error {
	return &StackTraceError{
		err:   errors.New(msg),
		stack: captureStack(3, defaultStackDepth), // Skip New + runtime frames
	}
}

// Newf creates a new error with stack trace.
func Newf(format string, args ...any) error {
	return &StackTraceError{
		err:   fmt.Errorf(format, args...),
		stack: captureStack(3, defaultStackDepth), // Skip New + runtime frames
	}
}

// Wrap wraps an existing error and adds a stack trace.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	// Avoid double-wrapping with a new stack trace
	if _, ok := errors.AsType[*StackTraceError](err); ok {
		return fmt.Errorf("%s: %w", msg, err)
	}

	return &StackTraceError{
		err:   fmt.Errorf("%s: %w", msg, err),
		stack: captureStack(3, defaultStackDepth),
	}
}

func formatStack(pcs []uintptr) []string {
	frames := runtime.CallersFrames(pcs)
	var trace []string

	for {
		frame, more := frames.Next()
		if !(filterOutPackageOn && strings.Contains(frame.File, "/"+filterOutPackage+"/")) {
			trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}
	return trace
}

// captureStack returns a slice of stack frames.
// Raw memory addresses pointing to each function call in the call stack.
func captureStack(skip int, depth int) []uintptr {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs)
	return pcs[:n]
}
