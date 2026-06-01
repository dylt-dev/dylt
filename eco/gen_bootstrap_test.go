package eco

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

//go:embed content/*
var content embed.FS

func TestGenBootstrap(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	ctx.Mute()

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

	// Confirm ECO_GENNER_FILENAME_PREFIX is properly set
	envGennerFilenamePrefix, is := os.LookupEnv("ECO_GENNER_FILENAME_PREFIX")
	require.True(t, is)

	// Confirm ECO_GENNER_TESTNAME_PREFIX is properly set
	envGennerTestNamePrefix, is := os.LookupEnv("ECO_GENNER_TESTNAME_PREFIX")
	require.True(t, is)

	r := rand.NewSource(time.Now().UTC().UnixNano())

	// Generate tests
	t.Logf("Generate %d test(s) of depth=%d", nTests, depth)
	for i := range nTests {
		bbDecl := bytes.Buffer{}
		common.WriteDeclaration(ctx, depth, r, &bbDecl)
		decl := strings.TrimSpace(bbDecl.String())
		testName := fmt.Sprintf("%s%d", envGennerTestNamePrefix, i)
		// 	// execute template
		data := map[string]any{
			"depth":      depth,
			"testName":   testName,
			"testNumber": i,
			"typeDecl":   decl,
		}
		buf, err := content.ReadFile("content/DeepTestGenner.tmpl")
		require.NoError(t, err)
		require.NotNil(t, buf)
		tmpl, err := template.New("genTestStage1").Parse(string(buf))

		filename := fmt.Sprintf("%s%d_test.go", envGennerFilenamePrefix, i)
		w, err := os.Create(filename)
		require.NoError(t, err)
		tmpl.Execute(w, data)
		err = w.Close()
		require.NoError(t, err)
	}
}
