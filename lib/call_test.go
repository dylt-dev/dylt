package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallDaylightScript (t *testing.T) {
	args := []string{"hello"}
	rc, err := CallDaylightScript(args)
	if err != nil { t.Fatalf("%v\n", err) }
	t.Logf("rc=%d\n", rc)
}


func TestCallDaylightScriptO (t *testing.T) {
	args := []string{"hello"}
	rc, stdout, err := CallDaylightScriptO(args)
	fmt.Printf("stdout=%s\n", stdout)
	if err != nil { t.Fatalf("%v\n", err) }
	t.Logf("rc=%d\n", rc)
}

func TestIsExecutable0 (t *testing.T) {
	path := "/tmp/i-am-program"
	flag, err := IsPathExecutable(path)
	assert.True(t, flag)
	assert.Nil(t, err)
}

func TestIsExecutable1 (t *testing.T) {
	path := "/tmp/i-am-file"
	flag, err := IsPathExecutable(path)
	assert.False(t, flag)
	assert.Nil(t, err)
}

func TestIsExecutable2 (t *testing.T) {
	path := "/tmp/i-am-not"
	flag, err := IsPathExecutable(path)
	assert.False(t, flag)
	assert.NotNil(t, err)
}

func TestRun0 (t *testing.T) {
	path := "/tmp/i-am-program"
	rc, stdout, err := RunScript(path, []string{})
	assert.Equal(t, 0, rc)
	assert.Equal(t, "hiii\n", string(stdout))
	assert.Nil(t, err)
}