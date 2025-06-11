package api

import (
	"io/fs"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	common.LocalInitLogging()
}

func TestAddContentFolder(t *testing.T) {
	fsContent, err := fs.Sub(EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmpl := template.New("content")
	_, err = tmpl.AddContentFS(fsContent)
	assert.NoError(t, err)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}
}

func TestAddContentFS(t *testing.T) {
	fsContent, err := fs.Sub(EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmpl := template.New("content")
	_, err = tmpl.AddContentFS(fsContent)
	assert.NoError(t, err)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}
}

func TestGetServiceTemplate(t *testing.T) {
	fsSvcFiles, err := fs.Sub(EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	svcName := "watch-daylight"
	tmpl, err := template.NewTemplate(fsSvcFiles, svcName)
	assert.NoError(t, err)
	require.NotNil(t, tmpl)
	for _, tmplChild := range tmpl.Templates() {
		t.Log(tmplChild.Name())
	}
}
