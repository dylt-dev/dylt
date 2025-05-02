package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunGet (t *testing.T) {
	key := "hello"
	err := RunGet(key)
	assert.Nil(t, err)
}

func TestGetCmd0 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	sCmdline := fmt.Sprintf("%s get hello", dyltPath)
	err := CheckRunCommandSuccess(sCmdline, t)
	assert.Nil(t, err)
}
