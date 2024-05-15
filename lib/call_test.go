package lib

import (
	"runtime/debug"
	"testing"
)

func TestCallDaylightScript (t *testing.T) {
	args := []string{}
	rc, err := CallDaylightScript(args)
	if err != nil { t.Fatalf("%v\n%s", err, string(debug.Stack())) }
	t.Logf("rc=%d\n", rc)
}
