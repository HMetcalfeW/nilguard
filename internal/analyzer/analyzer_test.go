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
// At this stage (Round 1), the Analyzer is still a stub, so this test is
// expected to FAIL once the // want markers are in place. In Round 2 we will
// implement the analyzer logic to make this test pass.
func TestNilguard(t *testing.T) {
	// analysistest.TestData locates the "testdata" directory relative to
	// the package containing this test file.
	testdata := analysistest.TestData()

	// We run the analyzer on both the "ok" and "bad" packages. analysistest
	// will compare the analyzer's diagnostics with the // want annotations.
	analysistest.Run(t, testdata, Analyzer, "ok", "bad")
}
