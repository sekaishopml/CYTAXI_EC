package logger

import (
	"context"
	"log/slog"
	"os"
)

type slogLogger struct {
	inner *slog.Logger
}

func NewSlog(level string, format string) Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: l}

	var h slog.Handler
	if format == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}

	return &slogLogger{inner: slog.New(h)}
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.inner.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.inner.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.inner.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.inner.Error(msg, args...)
}

func (l *slogLogger) Fatal(msg string, args ...any) {
	l.inner.Error(msg, args...)
	os.Exit(1)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{inner: l.inner.With(args...)}
}

type ctxKey string

const loggerKey ctxKey = "logger"

func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return NewSlog("info", "json")
	}
	return l
}
