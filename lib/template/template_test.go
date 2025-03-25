package template

import (
	"io/fs"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init () {
	lib.InitLogging()
}

func TestAddContentFolder (t *testing.T) {
	fsContent, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmpl := New("content")
	_, err = tmpl.AddContentFS(fsContent)
	assert.NoError(t, err)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}
}

func TestAddContentFS (t *testing.T) {
	fsContent, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmpl := New("content")
	_, err = tmpl.AddContentFS(fsContent)
	assert.NoError(t, err)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}
}

func TestGetServiceTemplate (t *testing.T) {
	svcName := "watch-daylight"
	tmpl, err := GetServiceTemplate(svcName)
	assert.NoError(t, err)
	require.NotNil(t, tmpl)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}	
}

