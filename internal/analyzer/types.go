package analyzer

import "go/token"

// pointerUseInfo tracks how a single pointer-typed identifier is used within
// a single function body.
//
// The analyzer records:
//   - the position of the first "use" (selector, method call, or star deref),
//   - whether we have seen at least one qualifying nil-check for this pointer
//     in the same function body.
type pointerUseInfo struct {
	// firstPos is the position of the first recorded use of this pointer in
	// the current function body. It is used as the position for any diagnostic
	// we emit about this pointer.
	firstPos token.Pos

	// hasCheck is true if we have observed at least one qualifying nil-check
	// anywhere in the current function body for this pointer. A qualifying
	// check is defined by the v1 policy in doc.go.
	hasCheck bool
}
