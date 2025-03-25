package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunList (t *testing.T) {
	err := RunList()
	assert.Nil(t, err)
}

func TestListCmd0 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s list", dyltPath)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}
