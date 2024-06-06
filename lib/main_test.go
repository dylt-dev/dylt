package lib

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func TestMain (m *testing.M) {
	logfile, err := os.OpenFile("log.txt", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	logger := slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)
	slog.Info("Testing logger config")
	code := m.Run()
	os.Exit(code)
}