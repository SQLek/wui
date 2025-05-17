package designer

import (
	"fmt"
	"go/ast"
)

// Represents user interface whole project.
//
// From UI design we are interested in one function declarat
type Project struct {
	// Whole ast tree of currently loaded project
	// may contain user generated code
	// and bisnes logic, so we should preserve it.
	astFile *ast.File

	// From UI standpoint, we only interested in one function.
	// It must contains `varName := wui.NewWindow()`
	// It will also contain other UI instructions.
	astFunc *ast.FuncDecl

	// All valid UI nodes (Node) found in the main UI function.
	nodes []*Node
}

// Nodes returns all valid UI nodes (Node) found in the main UI function.
func (p *Project) Nodes() []*Node {
	return p.nodes
}

// NewProject creates a Project from an ast.File, parsing all valid UI nodes in advance.
// Returns error if no node of kind containing "Window" is found.
func NewProject(astFile *ast.File) (*Project, error) {
	if astFile == nil {
		return nil, fmt.Errorf("astFile is nil")
	}

	var mainFunc *ast.FuncDecl
	var nodes []*Node

	// Find all functions and collect valid UI nodes
	for _, decl := range astFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		ast.Inspect(fn, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			node := NewNode(assign)
			if node == nil {
				return true
			}

			nodes = append(nodes, node)
			if node.kind == "Window" {
				mainFunc = fn
			}

			return true
		})
	}

	if mainFunc == nil {
		return nil, fmt.Errorf("no function with 'wui.NewWindow' found")
	}

	return &Project{
		astFile: astFile,
		astFunc: mainFunc,
		nodes:   nodes,
	}, nil
}
