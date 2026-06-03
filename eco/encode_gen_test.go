package eco

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

/*

Run a command line script to call TestBootstrapEncode test with Depth and Count

TestBootstrapEncode will generate a number of TestEncodeGenDepthM_N functions,
where M is the depth and N is the test number. These tests will be distributed
over a number of files, with a max of 20 test per file. @note maybe the batch
size should be an envvar.

Each TestEncodeGenXXX test will generate a sincgle TestEncodeDepthN test. These
tests will also be distrubuted over one or more files.

Templates
	EncodeGenTestFile - Preamble to a Gen Test File
	EncodeGenTest - A 'test' that will generate a real test
	EncodeTestFile - Preamble to a test file
	EncodeTest	An actual encode test

It's tricky to keep straight what all the pieces are, and what their job is.
Let's work backwards from the end state.

End state: One of more files of TestEncodeXXX() unit tests
TestEncodeXXX()                a test method that performs an encode test, by
                               staticly represetnting tests data, and calling
					           testEncode() to so the work
TestEncode.tmpl                a template that generates a TestEncodeXXX() function
testEncode(t, testData)        a function that performs the work of a TestEncode test
TestGenEncodeXXX()             a test that generates a TestEncodeXXX() test and writes
                               it to the proper test file, possibly creating it in the
						       process, by staticly representing gen test data, and
						       calling testGenEncode() to do the work
TestGenEncode.tmpl             a template that generate a TestGenEncode() function
testGenEncode(t, testData)     a function that performs the work of a
                               TestGenEncode() test
TestBootstrapEncode()          a test that validates some envvars and writes
                               TestGenEncode() tests to test files, possibly
							   creating them in the process
testBootstrapEncode(t, data)   a function that performs the work of testBootstrapEncode()

There are sources of complexity here
- the various bits of logic
- generating N tests and populating M files
- coordinating the whole workflow

Here's a divide-and-conquer approach to tackle all 3
- Solve the test file generation_population problem first, by starting with
  TestBootstrapEncode(). Verify the proper number of tests, with the proper
  names, are appearing in the proper files, which also are properly named and
  of the correct number
- Solve the bits-of-logic problem in reverse by focusing on geneating a proper
  encode_xxx_test.go file with a proper TestEncodeXXX() test, and work backwards
- Define the workflow the old fashioned way: with a single function of easy-to-read
  blue collar code

Let's get to it!
*/

// Create M tests of depth N
func TestEncodeBootstrap(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	ctx.Mute()

	envvars := []string{"ECOGEN", "ECOGEN_BATCHSIZE", "ECOGEN_COUNT", "ECOGEN_DEPTH", "ECOGEN_FILENAME_PREFIX", "ECOGEN_GENFILENAME_PREFIX", "ECOGEN_GENTESTNAME_PREFIX", "ECOGEN_TESTNAME_PREFIX"}
	common.EnvConfirmExists(t, envvars...)
	batchSize := common.GetEnvInt(t, "ECOGEN_BATCHSIZE")
	count := common.GetEnvInt(t, "ECOGEN_COUNT")
	depth := common.GetEnvInt(t, "ECOGEN_DEPTH")
	filenamePrefix := common.GetEnvString(t, "ECOGEN_GENFILENAME_PREFIX")
	testnamePrefix := common.GetEnvString(t, "ECOGEN_GENTESTNAME_PREFIX")

	testData := BootstrapEncodeData{
		BatchSize:      batchSize,
		Count:          count,
		Depth:          depth,
		FilenamePrefix: filenamePrefix,
		TestnamePrefix: testnamePrefix,
	}
	t.Log(testData)

	testBootstrapEncode(t, testData)
}

func testBootstrapEncode(t *testing.T, testData BootstrapEncodeData) {
	ctx := common.NewEcoContext(os.Stdout)
	ctx.Mute()

	// declarations
	decls := common.GenDeclarations(ctx, testData.Count, testData.Depth)

	// for 0 to Count
	for i, decl := range decls {
		// template data
		testname := testData.CreateTestName(i)
		data := EncodeGenTestData{
			Depth:      testData.Depth,
			Decl:       decl,
			TestName:   testname,
			TestNumber: i,
		}

		batchData := BatchData{
			BatchSize:      testData.BatchSize,
			FilenamePrefix: testData.FilenamePrefix,
			Total:          testData.Count,
		}

		// generate a new test, possibly in a new test file
		genEncodeGenTest(t, data, batchData)
	}
}

func newBatchData() BatchData {
	return BatchData{
		BatchSize:      20,
		FilenamePrefix: "testFile",
		Total:          1000,
	}
}

func testBatchDataFilename(t *testing.T, batchData BatchData, testNumber int, expectedBatchNumber int) {
	expected := fmt.Sprintf("%s_%02d_test.go", batchData.FilenamePrefix, expectedBatchNumber)
	filename := batchData.Filename(testNumber)
	t.Logf("filename=%s", filename)
	require.Equal(t, expected, filename)
}

func TestBatchDataFilename1(t *testing.T) {
	batchData := newBatchData()
	testNumber := 0
	expectedBatchNumber := 1
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename2(t *testing.T) {
	batchData := newBatchData()
	testNumber := 19
	expectedBatchNumber := 1
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename3(t *testing.T) {
	batchData := newBatchData()
	testNumber := 20
	expectedBatchNumber := 2
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename4(t *testing.T) {
	batchData := newBatchData()
	testNumber := 500
	expectedBatchNumber := 26
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename5(t *testing.T) {
	batchData := newBatchData()
	testNumber := 999
	expectedBatchNumber := 50
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename6(t *testing.T) {
	defer func() { r := recover(); require.NotNil(t, r) }()

	batchData := newBatchData()
	testNumber := 1000
	expectedBatchNumber := 0 // N/A
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestBatchDataFilename7(t *testing.T) {
	defer func() { r := recover(); require.NotNil(t, r) }()

	batchData := newBatchData()
	testNumber := 2000
	expectedBatchNumber := 0 // N/A
	testBatchDataFilename(t, batchData, testNumber, expectedBatchNumber)
}

func TestOpenEncodeGenTestFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "encode_gen_*_test.go")
	require.NoError(t, err)
	filename := tmpfile.Name()
	t.Log(filename)
	err = os.Remove(filename)

	w := openEncodeGenTestFile(t, filename)
	require.NotNil(t, w)
	err = w.(io.Closer).Close()
	require.NoError(t, err)

	buf, err := os.ReadFile(filename)
	require.NoError(t, err)
	require.NotNil(t, buf)
	t.Logf("\n%s", string(buf))
}

// Generate a new test, possibly in a new file
func genEncodeGenTest(t *testing.T, data EncodeGenTestData, batchData BatchData) {
	// Create filename
	filename := batchData.Filename(data.TestNumber)

	// open test file or create a new one
	w := openEncodeGenTestFile(t, filename)
	newline(t, w)
	newline(t, w)

	// Execute template
	tmpl := common.LoadTemplate(t, "content/EncodeGenTest.tmpl")
	err := tmpl.Execute(w, data)
	require.NoError(t, err)
}


// Generate a new test, possibly in a new file
func genEncodeTest(t *testing.T, data EncodeTestData, batchData BatchData) {
	ctx := common.NewEcoContext(os.Stdout)
	
	// Create filename
	filename := batchData.Filename(data.TestNumber)
	ctx.Infof("filename=%s", filename)

	// open test file or create a new one
	w := openEncodeGenTestFile(t, filename)
	newline(t, w)
	newline(t, w)

	// Execute template
	tmpl := common.LoadTemplate(t, "content/EncodeTest.tmpl")
	err := tmpl.Execute(w, data)
	require.NoError(t, err)
}


func newline(t *testing.T, w io.Writer) {
	_, err := w.Write([]byte("\n"))
	require.NoError(t, err)
}

// func getEncodeGenTestFilename(filenamePrefix string, testNumber int, testCount int) string {
// 	if testNumber >= testCount {
// 		panic(fmt.Errorf("testNumber must be less then testCount (%d>=%d)", testNumber, testCount))
// 	}
// 	batchSize := 20
// 	batchNum := testNumber / batchSize
// 	batchCount := testCount / batchSize

// 	nDigits := common.GetNumDigits(batchCount)
// 	fmtString := fmt.Sprintf("%s%%0%dd_test.go", filenamePrefix, nDigits)
// 	filename := fmt.Sprintf(fmtString, batchNum+1)

// 	return filename
// }

func openEncodeGenTestFile(t *testing.T, filename string) io.Writer {
	var w io.Writer
	w, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		if !os.IsNotExist(err) {
			require.NoError(t, err)
		}
		w, err = os.Create(filename)
		require.NoError(t, err)
		tmplNewFile := common.LoadTemplate(t, "content/EncodeGenTestFile.tmpl")
		err := tmplNewFile.Execute(w, nil)
		require.NoError(t, err)
	}

	return w
}

type EncodeTestData struct {
	Decl string
	Key KeyString
	ObjCtorStatements []string
	TestName string
	TestNumber int
}

func testEncodeGen(t *testing.T, rt reflect.Type, testNumber int) {
	ctx := common.NewEcoContext(os.Stdout)

	decl := common.GetDeclFromType(rt)
	values := common.GenScalarValues(ctx, rt)
	stmts := common.GenObjCtorStmts(ctx, rt, values)

	key := keyFromValues(values)
	
	testNamePrefix := common.GetEnvString(t, "ECOGEN_TESTNAME_PREFIX")
	testName := fmt.Sprintf("%s_%d", testNamePrefix, testNumber+1)
	// Pass the statements to the template
	data := EncodeTestData {
		Decl: decl,
		Key: key,
		ObjCtorStatements: stmts,
		TestName: testName,
		TestNumber: testNumber,
	}

	batchData := BatchData{
		BatchSize: common.GetEnvInt(t, "ECOGEN_BATCHSIZE"),
		FilenamePrefix: common.GetEnvString(t, "ECOGEN_FILENAME_PREFIX"),
		Total: common.GetEnvInt(t, "ECOGEN_COUNT"),
	}
	 
	genEncodeTest(t, data, batchData)
}
