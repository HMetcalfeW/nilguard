// Command nilguard-vet builds a vettool-compatible binary.
// Usage:
//
//	go build -o ./bin/nilguard-vet ./pkg/vettool
//	go vet -vettool=$(pwd)/bin/nilguard-vet ./...
package main

import (
	"github.com/HMetcalfe/nilguard/internal/analyzer"
	"github.com/HMetcalfe/nilguard/internal/runner"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	runner.SetupEnvDefaults()
	runner.RegisterEnvFlags()
	multichecker.Main(analyzer.Analyzer)
}
