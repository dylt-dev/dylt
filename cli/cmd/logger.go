package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"

	"github.com/dylt-dev/dylt/color"
)

type cliLogger struct {
	buf []byte
	*slog.Logger
}

func NewLogger(w io.Writer) *cliLogger {
	options := color.ColorOptions{Level: slog.LevelDebug}
	handler := color.NewColorHandler(w, options)
	return &cliLogger{
		Logger:  slog.New(handler),
		buf: make([]byte, 200),
	}
}

func (l *cliLogger) Append (s string) *cliLogger {
	l.buf = slices.Concat(l.buf, []byte(s))
	
	return l
}

func (l *cliLogger) Appendf (sfmt string, args ...any) *cliLogger {
	s := fmt.Sprintf(sfmt, args...)
	l.Append(s) 

	return l
}

func (l *cliLogger) AppendAndFlush (level slog.Level, s string) *cliLogger {
	l.Append(s)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)

	return l
}

func (l *cliLogger) AppendfAndFlush (level slog.Level, sfmt string, args ...any) {
	msg := fmt.Sprintf(sfmt, args...)
	l.Append(msg)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)
}

func (l *cliLogger) Debugf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Debug(l.indent() + s)
}

func (l *cliLogger) DebugContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.DebugContext(ctx, l.indent() + s)
}

func (l *cliLogger) Errorf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Error(l.indent() + s)
}

func (l *cliLogger) ErrorContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.ErrorContext(ctx, l.indent() + s)
}

func (l *cliLogger) Infof(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Info(l.indent() + s)
}

func (l *cliLogger) InfoContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.InfoContext(ctx, l.indent() + s)
}

func (l *cliLogger) Flush (level slog.Level) {
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.buf = make([]byte, 200)
}

func (l *cliLogger) Warnf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Warn(l.indent() + s)
}

func (l *cliLogger) WarnContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.WarnContext(ctx, l.indent() + s)
}

