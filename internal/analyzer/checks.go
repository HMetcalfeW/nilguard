package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

// baseIdentOf strips away syntactic noise like parentheses and returns the
// underlying *ast.Ident if the expression is a simple identifier (possibly
// wrapped in parentheses). If the base is not an identifier, it returns nil.
//
// Examples:
//
//	baseIdentOf(p)        -> ident "p"
//	baseIdentOf((*p).X)   -> nil (we only call this on the base expr)
//	baseIdentOf((p))      -> ident "p"
//	baseIdentOf(((*p)))   -> nil (underlying is a StarExpr, not an Ident)
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
//   - return
//   - panic(...)
//   - branch statements (break / continue / goto)
//
// This is intentionally conservative and coarse: treating break/continue/goto
// as exits simplifies the reasoning without affecting the core nil-check rule.
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
