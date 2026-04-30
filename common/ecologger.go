package common

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"

	"github.com/dylt-dev/dylt/color"
)

type Depther interface {
	Depth() int
}

type ecoLogger struct {
	buf []byte
	*slog.Logger
	depther Depther
}

func newEcoLogger(w io.Writer, depther Depther) *ecoLogger {
	options := color.ColorOptions{Level: slog.LevelDebug}
	handler := color.NewColorHandler(w, options)
	return &ecoLogger{
		Logger:  slog.New(handler),
		depther: depther,
		buf:     make([]byte, 200),
	}
}

func (l *ecoLogger) Append(s string) *ecoLogger {
	l.buf = slices.Concat(l.buf, []byte(s))

	return l
}

func (l *ecoLogger) Appendf(sfmt string, args ...any) *ecoLogger {
	s := fmt.Sprintf(sfmt, args...)
	return l.Append(s)
}

func (l *ecoLogger) AppendAndFlush(level slog.Level, s string) {
	l.Append(s)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)
}

func (l *ecoLogger) AppendfAndFlush(level slog.Level, sfmt string, args ...any) {
	msg := fmt.Sprintf(sfmt, args...)
	l.Append(msg)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)
}

func (l *ecoLogger) Debugf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Debug(l.Indent() + s)
}

func (l *ecoLogger) DebugContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.DebugContext(ctx, l.Indent()+s)
}

func (l *ecoLogger) Errorf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Error(l.Indent() + s)
}

func (l *ecoLogger) ErrorContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.ErrorContext(ctx, l.Indent()+s)
}

func (l *ecoLogger) Info(args ...any) {
	s := fmt.Sprint(args...)
	l.Logger.Info(l.Indent() + s)
}

func (l *ecoLogger) Infof(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Info(l.Indent() + s)
}

func (l *ecoLogger) InfoContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.InfoContext(ctx, l.Indent()+s)
}

func (l *ecoLogger) Flush(level slog.Level) {
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.buf = make([]byte, 200)
}

func (l *ecoLogger) Warnf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Warn(l.Indent() + s)
}

func (l *ecoLogger) WarnContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.WarnContext(ctx, l.Indent()+s)
}

func (l *ecoLogger) Comment(args ...any) {
	msg := fmt.Sprint(args...)
	l.Logger.Info(l.Indent() + string(color.Styledstring(msg).Fg(color.X11.CornflowerBlue)))
}

func (l *ecoLogger) Commentf(sFmt string, args ...any) {
	msg := fmt.Sprintf(sFmt, args...)
	l.Comment(msg)
}

func (l *ecoLogger) Indent() string {
	const tab = "  "
	return strings.Repeat(tab, l.depther.Depth())
}

// func (l *ecoLogger) Info(s string) {
// 	l.Logger.Info(l.Indent() + s)
// }

func (l *ecoLogger) Signature(name string, args ...any) {
	sig := CreateSignature(name, args...)
	l.Logger.Info(l.Indent() + sig)
}

// func CreateSignature(name string, args ...any) string {
// 	// highlight, concat, all that good stuff
// 	sFmt := fmt.Sprintf("%%s(%s)", strings.Repeat("%v, ", len(args)-1)+"%v")
// 	args2 := make([]any, len(args)+1)
// 	args2[0] = Highlight(name)
// 	for i, arg := range args {
// 		typ, is := arg.(reflect.Type)
// 		var sArg string
// 		if is {
// 			sArg = fmt.Sprintf("-%s-", FullTypeName(typ))
// 		} else {
// 			_, is := arg.(string)
// 			if is {
// 				sArg = fmt.Sprintf("\"%s\"", arg)
// 			} else {
// 				sArg = fmt.Sprintf("%v", arg)
// 			}
// 		}
// 		args2[i+1] = Lowlight(sArg)
// 	}
// 	s := fmt.Sprintf(sFmt, args2...)

// 	return s
// }
