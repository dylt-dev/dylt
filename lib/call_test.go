package lib

import (
	"fmt"
	"testing"
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

