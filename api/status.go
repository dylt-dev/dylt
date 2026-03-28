package api

import (
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)


func RunStatus() error {
	slog.Debug("RunStatus()")

	lib.RunStatus()

	return nil
}

