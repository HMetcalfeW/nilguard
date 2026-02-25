// Command nilguard runs the analyzer as a singlechecker CLI.
package main

import (
	"github.com/HMetcalfe/nilguard/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
