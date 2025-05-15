package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strconv"

	"github.com/gonutz/wui/v2"
)

func parseFile(window *wui.Window, path string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var foundFunc *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for function declarations
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true // Continue searching
		}

		// Check if this is a top-level function
		if funcDecl.Recv != nil {
			return true // Skip methods
		}

		// Look for wui.NewWindow() calls in the function body
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if it's a selector expression (like wui.NewWindow)
			sel, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Check if it's wui.NewWindow
			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "wui" && sel.Sel.Name == "NewWindow" {
					foundFunc = funcDecl
					return false // Stop searching
				}
			}
			return true
		})

		// If we found the function, stop searching
		if foundFunc != nil {
			return false
		}
		return true
	})

	if foundFunc != nil {
		return parseWindow(window, node, foundFunc)
	}

	return nil
}

type parsedElement struct {
	name        string
	node        node // node means node, not interface{}!
	font        *wui.Font
	methodCalls []*ast.CallExpr
}

func fontString(font *wui.Font) string {
	if font == nil {
		return ""
	}

	// Get font description
	desc := font.Desc

	// Build the string
	var result string
	result += "wui.NewFont(wui.FontDesc{\n"

	// Add font parameters
	if desc.Name != "" {
		result += fmt.Sprintf("\tName: \"%s\",\n", desc.Name)
	}
	if desc.Height != 0 {
		result += fmt.Sprintf("\tHeight: %d,\n", desc.Height)
	}
	if desc.Bold {
		result += "\tBold: true,\n"
	}
	if desc.Italic {
		result += "\tItalic: true,\n"
	}

	result += "})\n"
	return result
}

func (e parsedElement) String() string {
	if e.name == "" {
		return ""
	}

	if e.font != nil {
		return e.name + " := " + fontString(e.font)
	}

	// Build the string
	var result string
	result += fmt.Sprintf("%s := wui.New%s()\n", e.name, getElementType(e.node))

	// Add method calls
	for _, call := range e.methodCalls {
		sel := call.Fun.(*ast.SelectorExpr)
		args := make([]string, len(call.Args))
		for i, arg := range call.Args {
			args[i] = fmt.Sprintf("%v", arg)
		}
		result += fmt.Sprintf("\t.%s(%s)\n", sel.Sel.Name, args)
	}

	return result
}

func getElementType(n node) string {
	switch n.(type) {
	case *wui.Window:
		return "Window"
	case *wui.Button:
		return "Button"
	case *wui.Label:
		return "Label"
	case *wui.ComboBox:
		return "ComboBox"
	case *wui.ProgressBar:
		return "ProgressBar"
	case *wui.Slider:
		return "Slider"
	default:
		return "Unknown"
	}
}

func parseWindow(window *wui.Window, file *ast.File, funcDecl *ast.FuncDecl) error {
	elements := make(map[string]parsedElement)

	// First pass: find all assignments
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		if assignStmt.Tok != token.ASSIGN && assignStmt.Tok != token.DEFINE {
			return true
		}

		for _, expr := range assignStmt.Rhs {
			callExpr, ok := expr.(*ast.CallExpr)
			if !ok {
				continue
			}

			sel, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "wui" && len(sel.Sel.Name) > 3 && sel.Sel.Name[:3] == "New" {
					varName := assignStmt.Lhs[0].(*ast.Ident).Name
					element, ok := parseElement(callExpr)
					if !ok {
						continue
					}
					elements[varName] = element
					break
				}
			}
		}
		return true
	})

	// Second pass: find all method calls on our tracked variables
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Check if the receiver is one of our tracked variables
		if ident, ok := sel.X.(*ast.Ident); ok {
			if element, exists := elements[ident.Name]; exists {
				element.methodCalls = append(element.methodCalls, callExpr)
				elements[ident.Name] = element
			}
		}
		return true
	})

	fmt.Printf("Found %d elements:\n", len(elements))
	for _, element := range elements {
		fmt.Print(element.String())
	}

	return nil
}

func parseFont(callExpr *ast.CallExpr) (*wui.Font, error) {
	if len(callExpr.Args) != 1 {
		return nil, fmt.Errorf("NewFont requires exactly one argument")
	}

	// Try to parse the font description
	desc, ok := callExpr.Args[0].(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("NewFont argument must be a composite literal")
	}

	// Create font description
	fontDesc := wui.FontDesc{}
	for _, elt := range desc.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			slog.Warn("Font description must use key-value pairs")
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Name":
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				fontDesc.Name = lit.Value[1 : len(lit.Value)-1] // Remove quotes
			}
		case "Height":
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.INT {
				if height, err := strconv.Atoi(lit.Value); err == nil {
					fontDesc.Height = height
				}
			}
		case "Bold":
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.IDENT {
				fontDesc.Bold = lit.Value == "true"
			}
		case "Italic":
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.IDENT {
				fontDesc.Italic = lit.Value == "true"
			}
		}
	}

	return wui.NewFont(fontDesc)
}

func parseElement(callExpr *ast.CallExpr) (parsedElement, bool) {
	// Check if it's a selector expression (like wui.New...)
	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return parsedElement{}, false
	}

	// Check if it's from wui package
	if ident, ok := sel.X.(*ast.Ident); !ok || ident.Name != "wui" {
		return parsedElement{}, false
	}

	// Check if it's a New... function
	if len(sel.Sel.Name) <= 3 || sel.Sel.Name[:3] != "New" {
		return parsedElement{}, false
	}

	// Create the element based on the function name
	var element node
	var font *wui.Font

	switch sel.Sel.Name {
	case "NewWindow":
		element = wui.NewWindow()
	case "NewButton":
		element = wui.NewButton()
	case "NewLabel":
		element = wui.NewLabel()
	case "NewComboBox":
		element = wui.NewComboBox()
	case "NewProgressBar":
		element = wui.NewProgressBar()
	case "NewSlider":
		element = wui.NewSlider()
	case "NewFont":
		var err error
		font, err = parseFont(callExpr)
		if err != nil {
			slog.Warn("Failed to create font", "error", err)
			return parsedElement{}, false
		}
	default:
		slog.Warn("Unknown element", "name", sel.Sel.Name)
		return parsedElement{}, false
	}

	return parsedElement{
		name: sel.Sel.Name[3:], // Remove "New" prefix
		node: element,
		font: font,
	}, true
}
