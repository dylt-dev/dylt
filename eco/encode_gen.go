package eco

import (
	"fmt"

	"github.com/dylt-dev/dylt/common"
)

type BatchData struct {
	BatchSize      int
	FilenamePrefix string
	Total          int
}

func (d BatchData) Filename(testNumber int) string {
	fileNumber := d.batchNumber(testNumber)
	fmtString := fmt.Sprintf("%s_%%0%dd_test.go", d.FilenamePrefix, d.numDigits())
	filename := fmt.Sprintf(fmtString, fileNumber)
	return filename
}

func (d BatchData) batchNumber(testNumber int) int {
	// [0 .. batchSize-1] = 1, [batchSize .. 2*batchSize-1] = 2, etc
	if testNumber >= d.Total {
		panic(fmt.Errorf("testNumber must be less then total tests (%d>=%d)", testNumber, d.Total))
	}
	return testNumber/d.BatchSize + 1
}

func (d BatchData) numDigits() int {
	return common.GetNumDigits(d.Total / d.BatchSize)
}

type BootstrapEncodeData struct {
	BatchSize      int
	Count          int
	Depth          int
	FilenamePrefix string
	TestnamePrefix string
}

func (d BootstrapEncodeData) CreateTestName(testNumber int) string {
	testname := fmt.Sprintf("%s_%d", d.TestnamePrefix, testNumber+1)
	return testname
}

type EncodeGenTestData struct {
	Decl       string
	Depth      int
	TestName   string
	TestNumber int
}
