// +build linux

package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	memfd "github.com/justincormack/go-memfd"
)

func TestCreateMemFd (t *testing.T) {
	mfd, err := memfd.Create()
	assert.NoError(t, err)
	_, err = mfd.WriteString("hello")
}
