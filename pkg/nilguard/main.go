// Command nilguard runs the analyzer as a singlechecker CLI.
package main

import (
	"github.com/HMetcalfe/nilguard/internal/analyzer"
	"github.com/HMetcalfe/nilguard/internal/runner"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	runner.SetupEnvDefaults()
	runner.RegisterEnvFlags()
	singlechecker.Main(analyzer.Analyzer)
}
