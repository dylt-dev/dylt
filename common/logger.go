package common

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/dylt-dev/dylt/color"
)

type cliLogger struct {
	buf []byte
	*slog.Logger
}

func NewLogger(w io.Writer) *cliLogger {
	logLevel := getLogLevel()
	options := color.ColorOptions{Level: logLevel}
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

func (l *cliLogger) Comment(msg string) {
	l.Logger.Info(string(color.Styledstring(msg).Fg(color.X11.CornflowerBlue)))
}

func (l *cliLogger) Commentf(sFmt string, args ...any) {
	msg := fmt.Sprintf(sFmt, args...)
	l.Comment(msg)
}

func (l *cliLogger) Debugf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Debug(s)
}

func (l *cliLogger) DebugContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.DebugContext(ctx, s)
}

func (l *cliLogger) Errorf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Error(s)
}

func (l *cliLogger) ErrorContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.ErrorContext(ctx, s)
}

func (l *cliLogger) Infof(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Info(s)
}

func (l *cliLogger) InfoContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.InfoContext(ctx, s)
}

func (l *cliLogger) Flush (level slog.Level) {
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.buf = make([]byte, 200)
}

func (l *cliLogger) Signature(name string, args ...any) {
	sig := CreateSignature(name, args...)
	l.Logger.Debug(sig)
}

func CreateSignature(name string, args ...any) string {
	// highlight, concat, all that good stuff
	sFmt := fmt.Sprintf("%%s(%s)", strings.Repeat("%v, ", len(args)-1)+"%v")
	args2 := make([]any, len(args)+1)
	args2[0] = Highlight(name)
	for i, arg := range args {
		ty, is := arg.(reflect.Type)
		var sArg string
		if is {
			sArg = fmt.Sprintf("-%s-", FullTypeName(ty))
		} else {
			sArg = fmt.Sprintf("%v", arg)
		}
		args2[i+1] = Lowlight(sArg)
	}
	s := fmt.Sprintf(sFmt, args2...)

	return s
}

func (l *cliLogger) Warnf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Warn(s)
}

func (l *cliLogger) WarnContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.WarnContext(ctx, s)
}

func getLogLevel () slog.Leveler {
	envvar := os.Getenv("DEBUG")
	if envvar == "1" {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}