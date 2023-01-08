// Package logger offers the logging with structured log within a Google Cloud Logging.
// Example:
//
//	logger.Infof(request.Context(), "hello %s", "world")
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/brbranch/go-gae-logger/logger/model"
	"github.com/brbranch/go-gae-logger/logger/provider"
)

const DefaultCalller = 4

// CloseFunc must be called at the end of process.
type CloseFunc func()

type Logger struct {
	level  Level
	caller int
}

var lock sync.Mutex
var logger = &Logger{
	caller: DefaultCalller,
	level:  Debug,
}

// SetLevel set loglevel that should output logs.
func SetLevel(level Level) {
	lock.Lock()
	defer lock.Unlock()
	logger.level = level
}

// WithCaller returns Logger that set specified runtime.Caller number.
func WithCaller(caller int) *Logger {
	return &Logger{caller: caller, level: logger.level}
}

// WithoutCaller returns Logger that never logging source location.
func WithoutCaller() *Logger {
	return &Logger{caller: -1, level: logger.level}
}

func (l *Logger) out(ctx context.Context, std *os.File, level Level, format string, args ...interface{}) {
	if l.level > level {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	_ = json.NewEncoder(std).Encode(l.build(ctx, level, format, args...))
}

// Debugf records log at "DEBUG" severity.
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.out(ctx, os.Stdout, Debug, format, args...)
}

// Infof records log at "INFO" severity.
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.out(ctx, os.Stdout, Info, format, args...)
}

// Warnf records log at "WARNING" severity.
func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.out(ctx, os.Stdout, Warn, format, args...)
}

// Errorf records log at "ERROR" severity.
func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.out(ctx, os.Stderr, Error, format, args...)
}

// Fatalf records log at "CRITICAL" severity.
func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.out(ctx, os.Stderr, Fatal, format, args...)
}

// Span create Custom Span of Provider's Tracer.
func Span(ctx context.Context, label string) (context.Context, CloseFunc) {
	pv := provider.Get(ctx)
	if pv == nil {
		return ctx, func() {}
	}
	return pv.CustomSpan(ctx, label)
}

// Debugf calls global logger.Debugf.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args...)
}

// Infof calls global logger.Infof.
func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args...)
}

// Warnf calls global logger.Warnf.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(ctx, format, args...)
}

// Errorf calls global logger.Errorf.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(ctx, format, args...)
}

// Fatalf calls global logger.Fatalf.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logger.Fatalf(ctx, format, args...)
}

func (l *Logger) build(ctx context.Context, level Level, format string, args ...interface{}) *model.Payload {
	var source *model.SourceLocation = nil
	var spanId, traceId, projectId string
	if l.caller >= 0 {
		pc, file, line, ok := runtime.Caller(l.caller)
		if ok {
			source = &model.SourceLocation{
				File:     filepath.Base(file),
				Line:     line,
				Function: runtime.FuncForPC(pc).Name(),
			}
		}
	}
	pv := provider.Get(ctx)

	if pv != nil {
		span := pv.GetSpan(ctx)
		if span.Valid {
			spanId = span.SpanID
			traceId = fmt.Sprintf("projects/%s/traces/%s", projectId, traceId)
		}
		projectId = pv.ProjectID()

		return &model.Payload{
			Time:           time.Now(),
			SpanID:         spanId,
			Trace:          traceId,
			Message:        fmt.Sprintf(format, args...),
			Severity:       level.String(),
			SourceLocation: source,
		}
	}

	return &model.Payload{
		Time:           time.Now(),
		Message:        fmt.Sprintf(format, args...),
		Severity:       level.String(),
		SourceLocation: source,
	}
}
