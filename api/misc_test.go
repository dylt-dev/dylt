package api

import (
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestGenEtcdRunScript(t *testing.T) {
	type EtcdRunScriptData struct {
		Name                string
		DataDir             string
		AdvertiseClientUrls []string
		ListenClientUrls    []string
		ClientCertAuth      bool
	}
	data := EtcdRunScriptData{
		Name:                "arleytown",
		AdvertiseClientUrls: []string{"https:/127.0.0.1:2239", "ip2:2239"},
		ListenClientUrls:    []string{"https:/127.0.0.1:2239"},
		ClientCertAuth:      false,
	}
	buf, err := content.ReadFile("content/run-etcd.sh.tmpl")
	require.NoError(t, err)
	tmpl := template.New("hello")
	tmpl.Funcs(template.FuncMap{
		"join": strings.Join,
	})
	tmpl, err = tmpl.Parse(string(buf))
	require.NoError(t, err)
	err = tmpl.Execute(os.Stdout, data)
	require.NoError(t, err)
}

