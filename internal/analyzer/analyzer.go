package analyzer

import "golang.org/x/tools/go/analysis"

// Analyzer will be defined in the next step when we add real logic.
// Keeping a stub so the repo builds cleanly on day one.
var Analyzer = &analysis.Analyzer{
	Name: "nilguard",
	Doc:  "flags pointers used in a function without any nil check in that function (v1 policy)",
	Run:  func(*analysis.Pass) (interface{}, error) { return nil, nil },
}
