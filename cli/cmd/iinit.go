package cmd

import (
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/color"
)

var Logger *slog.Logger

func init () {
	options := color.ColorOptions{Level: slog.LevelDebug}
	handler := color.NewColorHandler(os.Stdout, options)
	Logger = slog.New(handler)
}