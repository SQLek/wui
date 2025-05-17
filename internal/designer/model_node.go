package designer

import "go/ast"

// Represent displatyable UI element.
//
// It will figure in code as a `varName := wui.NewSomethin()`
type Node struct {
	assign *ast.AssignStmt

	name, kind string
}

// NewNode creates a Node if assign is a valid assignment of the form
// varName := wui.NewSomething(). Returns nil if not valid.
func NewNode(assign *ast.AssignStmt) *Node {
	name, kind := wuiNameKind(assign)
	if name == "" || kind == "" {
		return nil
	}

	return &Node{
		assign: assign,
		name:   name,
		kind:   kind,
	}
}

func (n *Node) NameKind() (string, string) {
	return n.name, n.kind
}
