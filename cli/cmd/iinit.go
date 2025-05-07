package cmd

import (

	"log/slog"
	"github.com/dylt-dev/dylt/common"
)

var Logger *slog.Logger

func init () {
	options := common.ColorOptions{Level: slog.LevelDebug}
	handler := common.NewColorHandler(options)
	Logger = slog.New(handler)
}