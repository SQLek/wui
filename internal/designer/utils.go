package designer

import (
	"go/ast"
	"strings"
)

// wuiNameKind returns the name and kind of WUI element from the assignment statement.
// It checks if the assignment is of the form `varName := wui.NewSomething()`,
// and returns the variable name and the kind of WUI element.
// If the assignment is not valid, it returns empty strings.
func wuiNameKind(assign *ast.AssignStmt) (string, string) {
	if assign == nil || len(assign.Lhs) == 0 || len(assign.Rhs) == 0 {
		return "", ""
	}

	ident, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return "", ""
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return "", ""
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", ""
	}

	if !strings.HasPrefix(sel.Sel.Name, "New") {
		return "", ""
	}

	return ident.Name, sel.Sel.Name[3:]
}
