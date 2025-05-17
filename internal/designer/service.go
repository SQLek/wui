package designer

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ParseFile(fName string) (*Project, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, fName, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return NewProject(astFile)
}
