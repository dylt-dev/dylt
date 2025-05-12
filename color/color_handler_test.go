package color

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColorHandler0(t *testing.T) {
	options := ColorOptions{Level: slog.LevelDebug}
	handler := NewColorHandler(os.Stdout, options)
	logger := slog.New(handler)
	t.Logf("%#v", handler)
	logger.Debug("MEAT")
	logger.Info("hiii")
}

func TestColorHandler1(t *testing.T) {
	options := ColorOptions{Level: slog.LevelDebug}
	handler := NewColorHandler(os.Stdout, options)
	logger := slog.New(handler)
	t.Logf("Before: %v", logger.Handler().(*ColorHandler).meta)
	logger.WithGroup("g")
	logger.With(slog.String("foo", "13"))
	t.Logf("After: %v", logger.Handler().(*ColorHandler).meta)
	logger.Debug("test", "bar", "thirteen")
}

func TestCreateAttrMap (t *testing.T) {
	options := ColorOptions{Level: slog.LevelDebug}
	handler := NewColorHandler(os.Stdout, options)
	logger := slog.New(handler)
	logger.WithGroup("g")
	logger.With(slog.String("foo", "13"))

	attrMap, err := createAttrMap(logger.Handler().(*ColorHandler).meta)
	require.NoError(t, err)
	assert.NotEmpty(t, attrMap)
	t.Logf("%#v", attrMap)
}