// Exposes the Analyzer symbol for golangci-lint's custom linter plugin mechanism.
// Build with: GOFLAGS=-buildmode=plugin go build -o bin/nilguard.so ./plugin
package main

import "github.com/HMetcalfe/nilguard/internal/analyzer"

// Analyzer is exported so golangci-lint can discover and run it.
var Analyzer = analyzer.Analyzer
