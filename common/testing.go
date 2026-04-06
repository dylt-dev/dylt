package common

import (
	"testing"
)

func Setup (t *testing.T) func (t *testing.T) {
	t.Log("setup() ...")
	InitLogging()
	return Teardown
}

func Teardown (t *testing.T) {
	t.Log("teardown() ...")
}

