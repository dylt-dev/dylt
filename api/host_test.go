package api

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
)

func TestInstallWatchSvcService (t *testing.T) {
	if os.Getenv("RUNNING-ON-LINUX") != "Y" {
		t.Skip("Not running on a Linux system")
	}
	err := CreateWatchSvcService(2000, 2000)
	assert.NoError(t, err)
}

func TestRunHostInit(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	err := RunHostInit(501, 20)
	assert.Nil(t, err)
}

func TestHostInitCmd0(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s host init", dyltPath)
	err := lib.CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

// not-a-test
// Print out all the files in EMBED_SvcFiles. Useful sanity check.
func TestWalkSvcFolder(t *testing.T) {
	fs.WalkDir(EMBED_SvcFiles, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Printf("%s\n", p)
		}
		return nil
	})
}
