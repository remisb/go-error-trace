package trace

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

// filePaths holds the shared input data collected once in TestMain.
var filePaths []string

// samplePCs holds a realistic captured stack, collected once in TestMain.
var samplePCs []uintptr
var pHappy = "/usr/local/go/src/runtime/proc.go"
var pFail = "/Volumes/green/remis/go-error-trace/trace_bench_test.go"

// TestMain is the single entry point for all tests and benchmarks in this
// package. Setup runs once before any benchmark (or test) is executed.
func TestMain(m *testing.M) {
	filePaths = collectFilePaths()
	samplePCs = collectStack()
	os.Exit(m.Run())
}

// collectFilePaths resolves pcs into the file paths that formatStack would
// actually encounter, so the filter benchmarks operate on realistic data.
func collectFilePaths() []string {
	pcs := captureStack(2, defaultStackDepth)
	frames := runtime.CallersFrames(pcs)
	var paths []string
	for {
		frame, more := frames.Next()
		paths = append(paths, frame.File)
		if !more {
			break
		}
	}
	return paths
}

// collectStack captures a realistic stack to use as shared benchmark input.
func collectStack() []uintptr {
	pcs := make([]uintptr, defaultStackDepth)
	n := runtime.Callers(2, pcs)
	return pcs[:n]
}

// BenchmarkFilterHasSuffix measures the path.Dir + strings.HasSuffix approach
// used in the original formatStack.
func BenchmarkFilterHasSuffixHappyPath(b *testing.B) {
	b.ReportAllocs()
	_ = strings.HasSuffix(path.Dir(pHappy), filterOutPackage)
}

func BenchmarkFilterHasSuffixFailPath(b *testing.B) {
	b.ReportAllocs()
	_ = strings.HasSuffix(path.Dir(pFail), filterOutPackage)
}

func BenchmarkFormatStack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = formatStack(samplePCs)
	}
}
