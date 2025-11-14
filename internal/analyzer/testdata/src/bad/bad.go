package bad

// S is a sample struct used throughout the tests to model a pointer target.
type S struct {
	// X is a dummy field used for selector access in tests.
	X int
}

// M is a method on *S used to exercise method calls on pointer receivers.
func (s *S) M() {}

// noCheck demonstrates a simple violation: p is dereferenced via a selector
// in a function that never checks p for nil. This MUST produce a diagnostic.
func noCheck(p *S) {
	_ = p.X // want "pointer \"p\" is used in this function but never nil-checked"
}

// noCheckStar demonstrates a violation using a star dereference on *int.
func noCheckStar(p *int) {
	_ = *p // want "pointer \"p\" is used in this function but never nil-checked"
}

// nestedLiteralCheckOnly demonstrates that a check in a nested function
// literal does NOT satisfy uses in the outer function. The check inside
// the goroutine is intentionally ignored; the outer dereference must still
// produce a diagnostic.
func nestedLiteralCheckOnly(p *S) {
	go func() {
		if p != nil {
			_ = p.X
		}
	}()
	_ = p.X // want "pointer \"p\" is used in this function but never nil-checked"
}

// badEqualityCheck demonstrates that `if p == nil { ... }` without an early
// exit is not considered a qualifying check. The function continues after
// the if, so p must still be treated as unchecked for v1.
func badEqualityCheck(p *S) {
	if p == nil {
		p = &S{} // no early exit; this does not count as a guard in v1
	}
	_ = p.X // want "pointer \"p\" is used in this function but never nil-checked"
}

// multiPointers demonstrates that checks are tracked per-pointer identifier.
// p is checked, q is not. The use of q MUST produce a diagnostic.
func multiPointers(p, q *S) {
	if p != nil {
		_ = p.X
	}

	_ = q.X // want "pointer \"q\" is used in this function but never nil-checked"
}

// methodCall demonstrates a violation using a method call on a pointer receiver.
func methodCall(p *S) {
	p.M() // want "pointer \"p\" is used in this function but never nil-checked"
}
