package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallWatchSvcService (t *testing.T) {
	if os.Getenv("RUNNING-ON-LINUX") != "Y" {
		t.Skip("Not running on a Linux system")
	}
	err := CreateWatchSvcService(2000, 2000)
	assert.NoError(t, err)
}
