package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestNilguard runs the nilguard Analyzer against the packages under the
// local testdata directory. The testdata layout follows the convention:
//
//	testdata/src/ok   - packages where no diagnostics are expected
//	testdata/src/bad  - packages where diagnostics are annotated with // want
//
// The testdata layout uses // want markers to assert expected diagnostics.
func TestNilguard(t *testing.T) {
	// analysistest.TestData locates the "testdata" directory relative to
	// the package containing this test file.
	testdata := analysistest.TestData()

	// We run the analyzer on both the "ok" and "bad" packages. analysistest
	// will compare the analyzer's diagnostics with the // want annotations.
	analysistest.Run(t, testdata, Analyzer, "ok", "bad", "nolint")
}
