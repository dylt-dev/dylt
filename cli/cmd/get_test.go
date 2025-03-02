package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunGet (t *testing.T) {
	key := "hello"
	err := RunGet(key)
	assert.Nil(t, err)
}

func TestGetCmd0 (t *testing.T) {
	sCmdline := fmt.Sprintf("%s get hello", PATH_Dylt)
	err := CheckRunCommandSuccess(sCmdline, t)
	assert.Nil(t, err)
}
