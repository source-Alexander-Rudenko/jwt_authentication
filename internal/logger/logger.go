package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	Slog *slog.Logger
}

func NewTracer(l *slog.Logger, level tracelog.LogLevel) pgx.QueryTracer {
	return &tracelog.TraceLog{
		Logger:   &Logger{Slog: l},
		LogLevel: level,
	}
}

// Log вызывается PGX при событиях (start/end query и т.п.).
func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelError, tracelog.LogLevelNone:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	l.Slog.LogAttrs(ctx, lvl, msg, attrs...)
}

// Fatal Просто обертка для slog чтобы была эквивалентна Fatalf
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
