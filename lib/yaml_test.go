package lib

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYamlNode0(t *testing.T) {
	var node yaml.Node
	path := GetConfigFilePath()
	f, err := os.Open(path)
	assert.Nil(t, err)
	decoder := yaml.NewDecoder(f)
	decoder.Decode(&node)
	t.Logf("node.Tag=%s", node.Tag)
	t.Logf("node=%#v", node)
	for _, child := range node.Content {
		t.Logf("child.Tag=%s", child.Tag)
		t.Logf("child=%#v", child)
		for _, grandchild := range child.Content {
			t.Logf("grandchild.Tag=%s", grandchild.Tag)
			t.Logf("grandchild=%#v", grandchild)
		}
	}
}

func TestYamlMap1(t *testing.T) {
	sYaml := `
foo:
  bar: 13
`

	key := "foo.bar"
	f := strings.NewReader(sYaml)
	val, err := GetYamlValue(key, f)
	assert.NotNil(t, val)
	assert.Equal(t, 13, val.(int))
	assert.Nil(t, err)
}

func TestSetByKey(t *testing.T) {
	// Initial data
	data := map[string]any{}
	bar := map[string]any{"bar": 13}
	data["foo"] = bar
	// Set by key
	val := "mold"
	key := "yesterday.i.felt.so"
	SetKey(data, key, val)
	WriteYaml(data, os.Stdout)
}


func TestGetYamlValueNoKey (t *testing.T) {
	sYaml := `
foo:
  bar: 13
`
	f := strings.NewReader(sYaml)
	key := "XXX"
	val, err := GetYamlValue(key, f)
	assert.NotNil(t, err)
	assert.Nil(t, val)
}