package main

import (
	"bytes"
	"fmt"
	"go/format"
	"reflect"

	"github.com/gonutz/wui/v2"
)

func generateCode(w *wui.Window, isPreview bool) []byte {
	// TODO Remove the isPreview parameter once we can set window shortcuts
	// through the UI and generate them. Once we have that, temporarily add this
	// shortcut before generating the preview code and reset it afterwards, as
	// is done with the window position.
	var code bytes.Buffer
	code.WriteString(`package main

import "github.com/gonutz/wui/v2"

func main() {`)

	line := func(format string, a ...interface{}) {
		fmt.Fprint(&code, "\n")
		fmt.Fprintf(&code, format, a...)
	}

	name := names[w]
	if name == "" {
		name = defaultName(w)
	}
	writeControl(w, "", name, line)
	line("")
	if isPreview {
		line(name + ".SetShortcut(" + name + ".Close, wui.KeyEscape)")
	}
	line(name + ".Show()")
	code.WriteString("\n}")

	formatted, err := format.Source(code.Bytes())
	if err != nil {
		panic("We generated wrong code: " + err.Error())
	}
	return formatted
}

func writeControl(c interface{}, parentName, name string, line func(format string, a ...interface{})) {
	do := func(format string, a ...interface{}) {
		line(name+format, a...)
	}

	var fontName string
	if f, ok := c.(fonter); ok {
		font := f.Font()
		if font != nil {
			fontName = name + "Font"
			line(fontName + ", _ := wui.NewFont(wui.FontDesc{")
			if font.Desc.Name != "" {
				line("Name: %q,", font.Desc.Name)
			}
			if font.Desc.Height != 0 {
				line("Height: %d,", font.Desc.Height)
			}
			if font.Desc.Bold {
				line("Bold: true,")
			}
			if font.Desc.Italic {
				line("Italic: true,")
			}
			if font.Desc.Underlined {
				line("Underlined: true,")
			}
			if font.Desc.StrikedOut {
				line("StrikedOut: true,")
			}
			line("})")
			line("")
		}
	}

	typeName := reflect.TypeOf(c).Elem().Name()
	do(" := wui.New%s()", typeName)

	if fontName != "" {
		do(".SetFont(%s)", fontName)
	}

	setters := generateProperties(name, c)
	for _, setter := range setters {
		line("\t" + setter)
	}
	if parentName != "" {
		line("%s.Add(%s)", parentName, name)
	}
	line("")

	// TODO Generate ALL events.
	if p, ok := c.(*wui.PaintBox); ok {
		onPaint := event{p, "OnPaint"}
		if events[onPaint] != "" {
			do(".SetOnPaint(%s)", events[onPaint])
		}
	}

	if con, ok := c.(wui.Container); ok {
		for _, child := range con.Children() {
			childName := names[child]
			if childName == "" {
				childName = defaultName(child)
			}
			writeControl(child, name, childName, line)
		}
	}
}
