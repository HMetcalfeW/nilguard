package ok

// S is a sample struct used throughout the tests to model a pointer target.
type S struct {
	// X is a dummy field used for selector access in tests.
	X int
}

// M is a method on *S used to exercise method calls on pointer receivers.
func (s *S) M() {}

// guardedReturn demonstrates the pattern:
//
//	if p == nil { return } followed by uses of p.
//
// This MUST NOT produce a diagnostic.
func guardedReturn(p *S) {
	if p == nil {
		return
	}
	_ = p.X
	p.M()
}

// guardedBlock demonstrates the pattern:
//
//	if p != nil { ... uses of p ... }.
//
// This MUST NOT produce a diagnostic.
func guardedBlock(p *S) {
	if p != nil {
		_ = p.X
		p.M()
	}
}

// checkAfterUse demonstrates that, for v1, the location of the check
// within the function does not matter. As long as there is at least one
// qualifying check for p somewhere in the function, all uses are satisfied.
//
// This function MUST NOT produce a diagnostic, even though the first use
// of p occurs before the check.
func checkAfterUse(p *S) {
	_ = p.X
	if p != nil {
		_ = p.X
	}
}

// passOnly demonstrates that simply passing a pointer as an argument,
// without dereferencing or selecting from it, is NOT considered a "use"
// for v1. This function MUST NOT produce a diagnostic.
func passOnly(p *S) {
	consume(p)
}

// consume is a helper used to exercise passing pointers as arguments.
// It should not influence nilguard's behavior.
func consume(_ *S) {}

// nestedLiteralGuarded demonstrates that checks and uses entirely within a
// nested function literal do not affect the outer function.
//
// The outer function does not dereference p, so it MUST NOT produce a
// diagnostic, regardless of what happens inside the goroutine.
func nestedLiteralGuarded(p *S) {
	go func() {
		if p == nil {
			return
		}
		_ = p.X
	}()
}
