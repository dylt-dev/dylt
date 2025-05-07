package common

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/color"
)

type ColorHandler struct {
	options ColorOptions
}

func NewColorHandler(options ColorOptions) *ColorHandler {
	return &ColorHandler{options}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.options.Level == nil {
		return true
	}

	return level >= h.options.Level.Level()
}

func (h *ColorHandler) Handle(ctx context.Context, rec slog.Record) error {
	var s = color.Styledstring(rec.Message)
	if rec.Level == slog.LevelDebug {
		s = s.Fg(color.X11.gray30)
	}
	fmt.Println(s)

	return nil
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorHandler) WithGroup(grp string) slog.Handler {
	return h
}

type ColorOptions struct {
	Level slog.Leveler
}
