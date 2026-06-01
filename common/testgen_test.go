package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

type EncodeGennerData struct {
	Decl string
	Depth int
	Testname string
	TestNumber int
}

// Create M tests of depth N
func TestBootstrapEncode (t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	ctx.Mute()


	envvars := []string{"ECOGEN", "ECOGEN_COUNT", "ECOGEN_DEPTH", "ECOGEN_FILENAME_PREFIX", "ECOGEN_TESTNAME_PREFIX"}
	envConfirmExists(t, envvars...)	
	count := getEnvInt(t, "ECOGEN_COUNT")
	depth := getEnvInt(t, "ECOGEN_DEPTH")
	filenamePrefix := getEnvString(t, "ECOGEN_FILENAME_PREFIX")
	testnamePrefix := getEnvString(t, "ECOGEN_TESTNAME_PREFIX")

	// declaration
	decls := GenDeclarations(ctx, count, depth)

	// for 0 to Count
	for i, decl := range decls { 
	// template data
		testNumber := i+1
		testname := fmt.Sprintf("%s%d_%d", testnamePrefix, depth, testNumber)
		data := EncodeGennerData{
			Depth: depth,
			Decl: decl,
			Testname: testname,
			TestNumber: testNumber,
		}
		genEncodeGenner(t, data, i, count, filenamePrefix)
	}
}


/*
	type typ [][]struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}
	
	x0 := "meat"
	x1 := []string{x0}
	x2 := map[bool][]string{true: x1}
	x3 := []map[bool][]string{x2}
	x4 := map[int][]map[bool][]string{13: x3}
	x5 := struct{Dolorem map[int][]map[bool][]string}{Dolorem: x4}
	x6 := map[string]struct{Dolorem map[int][]map[bool][]string}{"foo": x5}
	x7 := struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}{Eum: x6}
	x8 := struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{Tempora: x7}
	x9 := []struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{x8}
	x10 := [][]struct{Tempora struct{Eum map[string]struct{Dolorem map[int][]map[bool][]string}}}{x9}
	var x typ = x10

	expected, err := json.Marshal(x0)
	require.NoError(t, err)
	kvs := encode(ctx, x)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, KeyString("/0/0/Tempora/Eum/foo/Dolorem/13/0/true/0"), kvs[0].Key)
	require.Equal(t, expected, kvs[0].Value)
	fmt.Fprint(t.Output(), kvs)
*/
func TestGenEncodeTest (t *testing.T) {

}


// Confirm envvars are set
func envConfirmExists (t *testing.T, names... string) {
	missing := []string{}
	
	for _, name := range names {
		_, is := os.LookupEnv(name)
		if !is {
			missing = append(missing, name)
		}
	}
	t.Skipf("Missing one or more require envvars (%s)", missing)
}


func genEncodeGenner (t *testing.T, data EncodeGennerData, i int, n int, filenamePrefix string) {
	// Create filename
	fileNumber := (i+1)/20
	nDigits := getNumDigits(n)
	fmtString := fmt.Sprintf("%s%%0%dd_test.go", filenamePrefix, nDigits)
	filename := fmt.Sprintf(fmtString, fileNumber)

	// open test file or create a new one
	var w io.Writer
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		if !os.IsNotExist(err) {
			require.NoError(t, err)
		}
		w, err = os.Create(filename)
		require.NoError(t, err)
		tmplNewFile := loadTemplate(t, "content/NewEncodeTestFile.tmpl")
		err := tmplNewFile.Execute(w, nil)
		require.NoError(t, err)
	}	
	tmpl := loadTemplate(t, "content/EncodeTestGen.tmpl")
	tmpl.Execute(w, data)
}


func getEnvInt (t *testing.T, name string) int {
	val, is := os.LookupEnv(name)
	require.True(t, is)
	n, err := strconv.Atoi(val)
	require.NoError(t, err)
	return n
}


func getEnvString (t *testing.T, name string) string {
	val, is := os.LookupEnv(name)
	require.True(t, is)
	return val
}


func getNumDigits (n int) int {
    if n == 0 {
        return 1
    }
    count := 0
    for n != 0 {
        n /= 10
        count++
    }
    
	return count
}


func loadTemplate (t *testing.T, tmplPath string) *template.Template {
	//	load template
	buf, err := content.ReadFile(tmplPath)
	require.NoError(t, err)
	require.NotNil(t, buf)
	tmplName := filepath.Base(tmplPath)
	tmpl, err := template.New(tmplName).Parse(string(buf))
	return tmpl
}