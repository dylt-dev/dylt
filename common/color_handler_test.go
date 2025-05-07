package common

import (
	"log/slog"
	"testing"
)

func TestColorHandler0(t *testing.T) {
	options := ColorOptions{Level: slog.LevelDebug}
	handler := NewColorHandler(options)
	logger := slog.New(handler)
	logger.Debug("MEAT")
	logger.Info("hiii")
}
