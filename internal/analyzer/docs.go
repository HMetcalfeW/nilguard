// Package analyzer implements the nilguard static analysis pass.
//
// v1 Policy (coarse, per-function):
//
//	If a pointer-typed identifier is used anywhere in a function body,
//	the same identifier must have at least one nil-check somewhere in that
//	function body. Nested functions (func literals) are excluded.
//
// Acceptable checks in v1:
//   - if p != nil { ... }
//   - if p == nil { return | panic(...) }  // early exit in 'then' branch
//
// Out of scope for v1:
//   - alias tracking (q := p)
//   - dominance / per-use flow
//   - interprocedural analysis
//
// Integration:
//   - Standalone CLI (singlechecker)
//   - go vet tool (multichecker)
//   - golangci-lint plugin (exported Analyzer symbol)
package analyzer
