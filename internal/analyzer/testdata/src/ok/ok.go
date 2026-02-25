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

// compoundAndCheck demonstrates that nil-checks within compound AND
// conditions are recognized:
//
//	if p != nil && p.X > 0 { ... }
//
// This MUST NOT produce a diagnostic.
func compoundAndCheck(p *S) {
	if p != nil && p.X > 0 {
		_ = p.X
	}
}

// compoundAndMultiple demonstrates that multiple pointers checked in a
// single compound AND condition are all recognized.
func compoundAndMultiple(p, q *S) {
	if p != nil && q != nil {
		_ = p.X
		_ = q.X
	}
}

// compoundOrEarlyReturn demonstrates that nil-checks within compound OR
// conditions followed by an early exit are recognized:
//
//	if p == nil || q == nil { return }
//
// This MUST NOT produce a diagnostic.
func compoundOrEarlyReturn(p, q *S) {
	if p == nil || q == nil {
		return
	}
	_ = p.X
	_ = q.X
}

// I is an interface for type assertion tests.
type I interface{ Foo() }

// Foo satisfies the I interface for *S.
func (s *S) Foo() {}

// typeAssertOk demonstrates that a two-value type assertion (v, ok := x.(*T))
// marks v as nil-checked, since the ok value guards it.
func typeAssertOk(x I) {
	v, ok := x.(*S)
	if ok {
		_ = v.X
	}
}

// typeSwitchCase demonstrates that a type switch (switch v := x.(type))
// marks v as nil-checked within each case clause.
func typeSwitchCase(x I) {
	switch v := x.(type) {
	case *S:
		_ = v.X
		v.M()
	}
}

// pointerReceiver demonstrates that a method with a pointer receiver
// can nil-check and use the receiver without diagnostic.
func (s *S) PointerReceiverGuarded() int {
	if s == nil {
		return 0
	}
	return s.X
}

// multiReturn demonstrates that a pointer obtained from a multi-return
// function is properly tracked. The nil-check satisfies the policy.
func multiReturn() {
	p, _ := getPointer()
	if p != nil {
		_ = p.X
	}
}

func getPointer() (*S, error) { return nil, nil }

// structFieldPointer demonstrates that pointers stored in struct fields
// are tracked when accessed via a local variable.
type Container struct {
	Ptr *S
}

func structFieldGuarded(c Container) {
	p := c.Ptr
	if p != nil {
		_ = p.X
	}
}
