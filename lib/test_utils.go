package lib

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func CheckRunCommandSuccess(sCmdlineArgs string, t *testing.T) error {
	dyltPath := GetAndValidateDyltPath(t)
	rc, stdout, err := RunCommand(dyltPath, strings.Fields(sCmdlineArgs)...)
	require.Equal(t, 0, rc)
	require.NotEmpty(t, stdout)
	require.Nil(t, err)
	t.Log(string(stdout))
	return err
}

func GetAndValidateDyltPath(t *testing.T) string {
	envName := "DYLT_PATH"
	dyltPath, is := os.LookupEnv(envName)
	if !is {
		t.Skipf("%s not set", envName)
	}
	_, err := os.Stat(dyltPath)
	if !os.IsNotExist(err) {
		t.Skipf("dylt path not found: %s", dyltPath)
	}
	return dyltPath
}

