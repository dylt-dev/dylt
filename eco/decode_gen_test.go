package eco

import (
	"bytes"
	"embed"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

//go:embed content/*
var content embed.FS

func TestGenBootstrap (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	ctx.Mute = true
	
	// Confirm ECO_GEN_TESTS is set
	envGenTests, is := os.LookupEnv("ECO_GEN_TESTS")
	if !is || (envGenTests != "1" && strings.ToLower(envGenTests) != "y") {
		t.Skipf("%s not set or not set to 1 or Y/y", "ECO_GEN_TESTS")
	}

	// Confirm ECO_DEPTH is properly set
	envDepth, is := os.LookupEnv("ECO_DEPTH")
	require.True(t, is)
	depth, err := strconv.Atoi(envDepth)
	require.NoError(t, err)
	require.Greater(t, depth, 0)

	// Confirm ECO_TEST_COUNT is properly set
	envTestCount, is := os.LookupEnv("ECO_TEST_COUNT")
	require.True(t, is)
	nTests, err := strconv.Atoi(envTestCount)
	require.NoError(t, err)
	require.Greater(t, nTests, 0)

	// Confirm ECO_OUTPUT_FILE is properly set
	envOutputFile, is := os.LookupEnv("ECO_OUTPUT_FILE")
	require.True(t, is)
	w, err := os.OpenFile(envOutputFile, os.O_CREATE | os.O_WRONLY, os.ModePerm)
	require.NoError(t, err)

	r := rand.NewSource(time.Now().UTC().UnixNano())

	// Generate tests
	t.Logf("Generate %d test(s) of depth=%d", nTests, depth)
	sDecls := make([]string, nTests)
	for i := range nTests {
		bbDecl := bytes.Buffer{}
		common.GenDeclaration(ctx, depth, r, &bbDecl)
		sDecls[i] = strings.TrimSpace(bbDecl.String())
	}

	// execute template
	data := map[string]any {
		"depth": depth,
		"typeDeclarations": sDecls,

	}
	buf, err := content.ReadFile("content/TestGenTest.tmpl")
	require.NoError(t, err)
	require.NotNil(t, buf)
	tmpl, err := template.New("genTestStage1").Parse(string(buf))
	tmpl.Execute(w, data)
	err = w.Close()
	require.NoError(t, err)
	// t.Log(sDecls)
}
