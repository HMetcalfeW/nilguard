// Command nilguard-vet builds a vettool-compatible binary.
// Usage:
//
//	go build -o ./bin/nilguard-vet ./cmd/vettool
//	go vet -vettool=$(pwd)/bin/nilguard-vet ./...
package main

import (
	"github.com/HMetcalfe/nilguard/internal/analyzer"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(analyzer.Analyzer)
}
