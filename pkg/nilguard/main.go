// Command nilguard runs the analyzer as a singlechecker CLI.
// Implementation wired after analyzer code lands.
package main

import (
	"github.com/HMetcalfe/nilguard/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
