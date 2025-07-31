package common

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLogLevelDebug (t *testing.T) {
	err := os.Setenv("DEBUG", "1")
	require.NoError(t, err)
	logLevel := getLogLevel()
	require.Equal(t, slog.LevelDebug, logLevel)
}

func TestGetLogLevelInfo0 (t *testing.T) {
	var err error
	err = os.Setenv("DEBUG", "0")
	require.NoError(t, err)
	logLevel := getLogLevel()
	require.Equal(t, slog.LevelInfo, logLevel)

}

func TestGetLogLevelInfo1 (t *testing.T) {
	var err error
	err = os.Setenv("DEBUG", "1")
	require.NoError(t, err)
	err = os.Unsetenv("DEBUG")
	require.NoError(t, err)
	logLevel := getLogLevel()
	require.Equal(t, slog.LevelInfo, logLevel)

}

