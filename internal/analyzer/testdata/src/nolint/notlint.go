package nolint

// S is a sample struct mirroring the other test packages.
type S struct {
	// X is a dummy field used for selector access in tests.
	X int
}

// M is a method on *S used to exercise method calls on pointer receivers.
func (s *S) M() {}

// lineSuppressionSelector demonstrates that a //nolint:nilguard comment on the
// same line as a pointer use suppresses the diagnostic.
func lineSuppressionSelector(p *S) {
	_ = p.X //nolint:nilguard
}

// lineSuppressionMethod demonstrates suppression for a method call.
func lineSuppressionMethod(p *S) {
	p.M() // nolint:nilguard
}

// multiLinterSuppression demonstrates suppression when nilguard appears in a
// comma-separated list of linters inside the nolint directive.
func multiLinterSuppression(p *S) {
	_ = p.X // nolint:foo,nilguard,bar
}

// noSuppressionOtherLinter demonstrates that a nolint directive for other
// linters, without mentioning nilguard, MUST NOT suppress our diagnostic.
//
// This function SHOULD trigger a diagnostic from nilguard.
func noSuppressionOtherLinter(p *S) {
	_ = p.X // nolint:foo,bar // want "pointer \"p\" is used in this function but never nil-checked"
}
