// Package analyzer implements the core nilguard static analysis pass.
//
// # v1 Policy (Per-Function, Coarse)
//
// For each function (a top-level function declaration or a function literal),
// nilguard enforces the following rule:
//
//	If a pointer-typed identifier is *used* anywhere in the function body,
//	that same identifier must have at least one qualifying nil-check somewhere
//	in that function body. Nested function literals are treated as separate
//	functions and do not share checks or uses with their enclosing functions.
//
// A "use" (for v1) is any of the following:
//
//   - Star dereference: *p
//   - Selector on a pointer base: p.Field
//   - Method call on a pointer base: p.Method()
//
// Parenthesized forms such as (*p).Field and (*p).Method() are conceptually
// treated the same as the unparenthesized forms.
//
// A "qualifying nil-check" (for v1) is any of:
//
//   - An if statement whose condition is `p != nil`.
//   - An if statement whose condition is `p == nil` and whose "then" branch
//     exits the function early via return or panic(...).
//
// Examples of checks that DO count:
//
//	if p != nil {
//	    // ...
//	}
//
//	if p == nil {
//	    return
//	}
//
//	if p == nil {
//	    panic("nil pointer")
//	}
//
// Examples of checks that do NOT count:
//
//	if p == nil {
//	    p = new(T) // no early exit; function continues
//	}
//
//	if p == nil || someOtherCond {
//	    // complex condition; out of scope for v1
//	}
//
// # Out of Scope for v1
//
// The following are intentionally out of scope for the initial implementation:
//
//   - Alias tracking: q := p; uses of q are not associated back to p.
//   - Interprocedural reasoning: constructors like NewT() are not treated
//     specially, even if they always return non-nil.
//   - Dominance / per-use flow: a single qualifying check anywhere in the
//     function satisfies all uses of the pointer in that function.
//   - Checks or uses inside nested function literals: a func literal is
//     treated as its own function for nilguard's purposes.
//
// # Integrations
//
// The Analyzer type defined in this package is designed to be reused in
// multiple frontends:
//
//   - Standalone CLI (via x/tools/go/analysis/singlechecker).
//   - go vet tool (via x/tools/go/analysis/multichecker).
//   - golangci-lint plugin (exported Analyzer symbol in a plugin package).
//
// The Analyzer itself does not perform any I/O beyond reporting diagnostics
// through analysis.Pass.
package analyzer
