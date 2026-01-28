package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the nilguard analysis pass entrypoint.
//
// It enforces the v1 policy described in doc.go:
//
//	For each function (FuncDecl or FuncLit), if a pointer-typed identifier is
//	used via dereference, selector, or method call anywhere in the function
//	body, that same identifier must have at least one qualifying nil-check
//	somewhere in that body. Nested function literals are treated as separate
//	functions and do not share state with their enclosing functions.
//
// The Analyzer does not perform any I/O beyond reporting diagnostics through
// the provided analysis.Pass.
var Analyzer = &analysis.Analyzer{
	Name: "nilguard",
	Doc:  "flags pointers used in a function without any nil check in that function (v1 policy)",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

var excludeTests bool

func init() {
	Analyzer.Flags.BoolVar(&excludeTests, "exclude-tests", false, "exclude _test.go files from analysis")
}

// run is the main entrypoint invoked by the analysis framework. It retrieves
// the precomputed inspector and applies our per-function analysis to each
// function declaration and function literal in the package.
func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Precompute an index of lines that have a nolint directive for nilguard.
	noLintIndex := buildNoLintIndex(pass)
	fileIndex := buildFileIndex(pass)

	// We care about function declarations and function literals. Both are
	// treated the same from the perspective of our rule: each function body
	// is analyzed independently.
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		if excludeTests && isTestFile(pass.Fset, n.Pos()) {
			return
		}
		var body *ast.BlockStmt

		switch fn := n.(type) {
		case *ast.FuncDecl:
			// Methods and plain functions both appear as FuncDecl. If there is
			// no body (e.g. an external declaration), there is nothing to do.
			if fn.Body == nil {
				return
			}
			body = fn.Body

		case *ast.FuncLit:
			body = fn.Body
		}

		checkFunc(pass, body, noLintIndex, fileIndex)
	})

	return nil, nil
}

// checkFunc performs the per-function analysis for a single function body.
//
// It walks the body, recording:
//   - pointer uses (selectors, method calls, star dereferences), and
//   - qualifying nil-checks (if p != nil, if p == nil { early-exit }).
//
// Nested function literals are skipped entirely: they are treated as separate
// functions by the outer run() traversal, and their checks/uses do not
// affect the enclosing function.
//
// At the end of the traversal, any pointer that was used at least once but
// never nil-checked will result in a single diagnostic at its first use.
func checkFunc(pass *analysis.Pass, body *ast.BlockStmt, noLintIndex map[*token.File]map[int]bool, fileIndex map[string]bool) {
	// ptrs maps each pointer-typed identifier (by its *ast.Object) to its
	// usage information within this function body.
	ptrs := make(map[types.Object]*pointerUseInfo)

	// recordUse registers a "use" of a pointer at the given position. A use
	// is any selector, method call, or star dereference whose base expression
	// is a pointer-typed identifier.
	recordUse := func(id *ast.Ident, pos token.Pos) {
		if id == nil {
			return
		}
		if !isPointerIdent(pass.TypesInfo, id) {
			return
		}

		// Look up the canonical types.Object behind this identifier.
		obj := pass.TypesInfo.ObjectOf(id)
		if obj == nil {
			// Without a types.Object, we cannot reliably track this identifier.
			return
		}

		if _, ok := ptrs[obj]; !ok {
			ptrs[obj] = &pointerUseInfo{
				firstPos: pos,
				hasCheck: false,
			}
		}
	}

	// markChecked notes that we have seen at least one qualifying nil-check
	// for the given pointer within this function body.
	markChecked := func(id *ast.Ident) {
		if id == nil {
			return
		}
		if !isPointerIdent(pass.TypesInfo, id) {
			return
		}

		obj := pass.TypesInfo.ObjectOf(id)
		if obj == nil {
			return
		}

		info, ok := ptrs[obj]
		if !ok {
			info = &pointerUseInfo{}
			ptrs[obj] = info
		}
		info.hasCheck = true
	}

	// Walk the function body. We explicitly skip nested function literals,
	// as those are analyzed separately as their own functions.
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			// Do not descend into nested function literals; they have their own
			// checkFunc invocation via the outer run() traversal.
			return false

		case *ast.StarExpr:
			// *p: record a use if the base is a pointer-typed identifier.
			if id := baseIdentOf(x.X); id != nil {
				recordUse(id, x.Pos())
			}

		case *ast.SelectorExpr:
			// p.Field or p.Method: record a use if the receiver/base is a
			// pointer-typed identifier. Parentheses around the base are
			// handled by baseIdentOf.
			if id := baseIdentOf(x.X); id != nil {
				recordUse(id, x.Pos())
			}

		case *ast.IfStmt:
			// if p != nil { ... } is always treated as a qualifying check.
			if id := binopPtrNil(pass.TypesInfo, x.Cond, token.NEQ); id != nil {
				markChecked(id)
			}

			// if p == nil { early-exit } is also treated as a qualifying check.
			if id := binopPtrNil(pass.TypesInfo, x.Cond, token.EQL); id != nil {
				if exitsEarly(x.Body) {
					markChecked(id)
				}
			}
		}

		return true
	})

	// Emit diagnostics for any pointer that was used but never nil-checked.
	for obj, info := range ptrs {
		if info.firstPos == 0 {
			// Pointer never used; nothing to report.
			continue
		}
		if info.hasCheck {
			// A qualifying nil-check exists; our v1 policy is satisfied.
			continue
		}

		// Skip diagnostics for files outside the current package's file set.
		if !isFileInPackage(pass.Fset, fileIndex, info.firstPos) {
			continue
		}

		// Respect per-line //nolint:nilguard directives.
		if hasNoLintNilguard(pass.Fset, noLintIndex, info.firstPos) {
			continue
		}

		// Report a single diagnostic per pointer at its first use position.
		pass.Reportf(info.firstPos,
			"pointer %q is used in this function but never nil-checked",
			obj.Name())
	}

}
