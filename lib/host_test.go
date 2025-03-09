package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallWatchSvcService (t *testing.T) {
	err := CreateWatchSvcService(2000, 2000)
	assert.NoError(t, err)
}
