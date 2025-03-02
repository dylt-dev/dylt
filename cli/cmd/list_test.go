package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunList (t *testing.T) {
	err := RunList()
	assert.Nil(t, err)
}

func TestListCmd0 (t *testing.T) {
	cmd := fmt.Sprintf("%s list", PATH_Dylt)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}
