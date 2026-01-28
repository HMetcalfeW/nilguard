package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// baseIdentOf strips away syntactic noise like parentheses and returns the
// underlying *ast.Ident if the expression is a simple identifier (possibly
// wrapped in parentheses). If the base is not an identifier, it returns nil.
//
// Examples:
//
//	baseIdentOf(p)      -> ident "p"
//	baseIdentOf((p))    -> ident "p"
//	baseIdentOf((*p).X) -> nil (we only call this on the base expr)
//
// Note: This helper is intentionally conservative and only recognizes simple
// identifiers as bases. More complex expressions are left to later versions.
func baseIdentOf(expr ast.Expr) *ast.Ident {
	for {
		switch e := expr.(type) {
		case *ast.Ident:
			return e
		case *ast.ParenExpr:
			expr = e.X
		default:
			return nil
		}
	}
}

// isPointerIdent reports whether id has a pointer underlying type according
// to the provided types.Info. If type information is missing, it returns false.
func isPointerIdent(info *types.Info, id *ast.Ident) bool {
	if id == nil {
		return false
	}
	t := info.TypeOf(id)
	if t == nil {
		return false
	}
	_, ok := t.Underlying().(*types.Pointer)
	return ok
}

// isNil reports whether e is the predeclared identifier "nil".
func isNil(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "nil"
}

// binopPtrNil inspects e and, if it is a binary expression of the form
//
//	p <op> nil   or   nil <op> p
//
// where <op> matches want and p is a pointer-typed identifier, returns the
// *ast.Ident for p. Otherwise it returns nil.
//
// This is used to recognize `p == nil` and `p != nil` conditions in if
// statements.
func binopPtrNil(info *types.Info, e ast.Expr, want token.Token) *ast.Ident {
	b, ok := e.(*ast.BinaryExpr)
	if !ok || b.Op != want {
		return nil
	}

	// Match: id <op> nil
	if id, ok := b.X.(*ast.Ident); ok && isPointerIdent(info, id) && isNil(b.Y) {
		return id
	}

	// Match: nil <op> id
	if id, ok := b.Y.(*ast.Ident); ok && isPointerIdent(info, id) && isNil(b.X) {
		return id
	}

	return nil
}

// exitsEarly reports whether the given block ends with an unconditional exit
// from the current function according to our v1 policy.
//
// For v1, we consider the following as "early exits":
//
//   - return
//   - panic(...)
//   - branch statements (break / continue / goto)
//
// This is intentionally conservative and coarse: treating break/continue/goto
// as exits simplifies the reasoning without affecting the core nil-check rule
// in real-world code.
func exitsEarly(b *ast.BlockStmt) bool {
	if b == nil || len(b.List) == 0 {
		return false
	}

	last := b.List[len(b.List)-1]

	switch s := last.(type) {
	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		// Consider any branch as an exit for v1. This may be conservative in
		// some cases but keeps the implementation simple.
		return s.Tok == token.GOTO || s.Tok == token.BREAK || s.Tok == token.CONTINUE

	case *ast.ExprStmt:
		// Match panic(...)
		if call, ok := s.X.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "panic" {
				return true
			}
		}
	}

	return false
}

// buildNoLintIndex constructs an index of lines that contain a nolint directive
// mentioning "nilguard", grouped by file.
//
// The result maps each *token.File to a set of line numbers (1-based) for
// which diagnostics from nilguard should be suppressed.
//
// We only consider comment texts that contain both "nolint" and "nilguard"
// (case-insensitive), which supports common forms such as:
//
//	//nolint:nilguard
//	// nolint:nilguard
//	// nolint:foo,nilguard,bar
//
// but does not treat bare //nolint as sufficient to suppress nilguard.
func buildNoLintIndex(pass *analysis.Pass) map[*token.File]map[int]bool {
	index := make(map[*token.File]map[int]bool)

	for _, f := range pass.Files {
		if f == nil {
			continue
		}

		tf := pass.Fset.File(f.Pos())
		if tf == nil {
			continue
		}

		// Lazily allocate the per-file map when we actually find a relevant comment.
		var lines map[int]bool

		for _, cg := range f.Comments {
			if cg == nil {
				continue
			}
			for _, c := range cg.List {
				if c == nil {
					continue
				}
				// Normalize the comment text for a simple substring check.
				text := strings.ToLower(c.Text)

				// Quick filter: require both "nolint" and "nilguard".
				if !strings.Contains(text, "nolint") || !strings.Contains(text, "nilguard") {
					continue
				}

				if lines == nil {
					lines = make(map[int]bool)
					index[tf] = lines
				}

				line := tf.Line(c.Slash)
				if line > 0 {
					lines[line] = true
				}
			}
		}
	}

	return index
}

// buildFileIndex records the set of file paths in the current package.
func buildFileIndex(pass *analysis.Pass) map[string]bool {
	index := make(map[string]bool)
	if pass == nil || pass.Fset == nil {
		return index
	}

	for _, f := range pass.Files {
		if f == nil {
			continue
		}
		tf := pass.Fset.File(f.Pos())
		if tf == nil {
			continue
		}
		index[tf.Name()] = true
	}

	return index
}

// isFileInPackage reports whether pos belongs to a file that appears in the
// current package's file set.
func isFileInPackage(fset *token.FileSet, index map[string]bool, pos token.Pos) bool {
	if fset == nil {
		return false
	}
	tf := fset.File(pos)
	if tf == nil {
		return false
	}
	return index[tf.Name()]
}

// hasNoLintNilguard reports whether the source line corresponding to pos in
// the given file set is marked as having a nolint directive for nilguard in
// the provided index.
//
// The index is expected to be produced by buildNoLintIndex.
func hasNoLintNilguard(fset *token.FileSet, index map[*token.File]map[int]bool, pos token.Pos) bool {
	if fset == nil {
		return false
	}

	tf := fset.File(pos)
	if tf == nil {
		return false
	}

	lines, ok := index[tf]
	if !ok {
		return false
	}

	line := tf.Line(pos)
	if line <= 0 {
		return false
	}

	return lines[line]
}

// isTestFile reports whether the file containing pos ends with _test.go.
func isTestFile(fset *token.FileSet, pos token.Pos) bool {
	if fset == nil {
		return false
	}
	tf := fset.File(pos)
	if tf == nil {
		return false
	}
	return strings.HasSuffix(tf.Name(), "_test.go")
}
