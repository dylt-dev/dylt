package cmd

import (
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/common"
)

var Logger *slog.Logger

func init () {
	options := common.ColorOptions{Level: slog.LevelDebug}
	handler := common.NewColorHandler(os.Stdout, options)
	Logger = slog.New(handler)
}