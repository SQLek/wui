package main

import (
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gonutz/w32/v2"
	"github.com/gonutz/wui/v2"
)

// TODO Have color property for window background.

// TODO Have icon for window.

// TODO Have cursor properties for all controls, first let all controls have
// changeable cursors.

// TODO Edit main menu.

// TODO Have a way to edit shortcuts.

// TODO Make edit lines select the whole text when they receive focus.

// TODO Un-highlight the template when the mouse leaves the palette area, even
// when it leaves it fast, in which case it does not receive a mouse move
// message.

// TODO Have a way to hide the app icon (WS_EX_DLGMODALFRAME).

// TODO Have a way to hide the border completely but make it still resizeable.

var (
	// names associates variable names with the controls.
	names = make(map[interface{}]string)

	// The build* variables are for previews which create temporary Go and exe
	// files that need to be deleted at the end of the program.
	buildDir    = "."
	buildPrefix = "wui_designer_preview_"
	buildCount  = 0

	events = make(map[event]string)
)

type event struct {
	control interface{}
	name    string
}

type Designer struct {
	w                *wui.Window
	active           node
	preview          *wui.PaintBox
	updateProperties func()
}

func main() {
	// Create a temporary directory to save our preview builds in.
	if dir, err := os.MkdirTemp("", "wui_designer_preview_builds"); err == nil {
		buildDir = dir
		defer os.Remove(dir)
	}
	// After closing the designer, delete all preview builds from this session.
	defer func() {
		if files, err := os.ReadDir(buildDir); err == nil {
			for _, file := range files {
				if !file.IsDir() &&
					strings.HasSuffix(file.Name(), ".exe") &&
					strings.HasPrefix(file.Name(), buildPrefix) {
					// Ignore the error on the remove, one of the builds might
					// still be running and cannot be removed. This is OK, most
					// of the files will get deleted.
					os.Remove(filepath.Join(buildDir, file.Name()))
				}
			}
			os.Remove(filepath.Join(buildDir, "go.mod"))
			os.Remove(filepath.Join(buildDir, "go.sum"))
		}
	}()

	var (
		// innerX and Y is the top-left corner of where theWindow's inner
		// rectangle is drawn, relative to the application window. This means we
		// can use innerX and Y in the application window's mouse events to find
		// the relative mouse position inside theWindow.
		innerX, innerY int
		// active is the highlighted control whose properties are shown in the
		// tool bar.
		//active node
		// TODO Move preview somewhere else.

	)

	theWindow := defaultWindow()
	names[theWindow] = "window"

	font, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11})
	bold, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11, Bold: true})
	italic, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11, Italic: true})
	underlined, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11, Underlined: true})
	strikedOut, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11, StrikedOut: true})

	w := wui.NewWindow()
	w.SetFont(font)
	w.SetTitle("wui Designer")
	w.SetBackground(wui.ColorButtonFace)
	w.SetInnerSize(800, 600)

	d := Designer{
		w:       w,
		preview: wui.NewPaintBox(),
	}

	menu := wui.NewMainMenu()
	fileMenu := wui.NewMenu("&File")
	editMenu := wui.NewMenu("&Edit")
	fileOpenMenu := wui.NewMenuString("&Open File...\tCtrl+O")
	fileSaveMenu := wui.NewMenuString("&Save File\tCtrl+S")
	fileSaveAsMenu := wui.NewMenuString("Save File &As...\tCtrl+Shift+S")
	previewMenu := wui.NewMenuString("&Run Preview\tCtrl+R")
	exitMenu := wui.NewMenuString("E&xit\tAlt+F4")
	undoMenu := wui.NewMenuString("&Undo\tCtrl+Z")
	redoMenu := wui.NewMenuString("&Redo\tCtrl+Shift+Z")
	deleteMenu := wui.NewMenuString("&Delete\tCtrl+Del")
	fileMenu.Add(fileOpenMenu)
	fileMenu.Add(fileSaveMenu)
	fileMenu.Add(fileSaveAsMenu)
	fileMenu.Add(wui.NewMenuSeparator())
	fileMenu.Add(previewMenu)
	fileMenu.Add(wui.NewMenuSeparator())
	fileMenu.Add(exitMenu)
	editMenu.Add(undoMenu)
	editMenu.Add(redoMenu)
	editMenu.Add(wui.NewMenuSeparator())
	editMenu.Add(deleteMenu)
	menu.Add(fileMenu)
	menu.Add(editMenu)
	w.SetMenu(menu)

	// TODO Doing this after the menu does not work.
	//w.SetInnerSize(800, 600)

	uiProps := []uiProp{
		d.stringProp("Title", "Title"),
		d.stringProp("Text", "Text"),
		d.enumProp("Window State", "State",
			"Normal", "Maximized", "Minimized",
		),
		d.boolProp("Min Button", "HasMinButton"),
		d.boolProp("Max Button", "HasMaxButton"),
		d.boolProp("Close Button", "HasCloseButton"),
		d.boolProp("Has Border", "HasBorder"),
		d.boolProp("Resizable", "Resizable"),
		d.intProp("Alpha", "Alpha", 0, 255),
		d.boolProp("Enabled", "Enabled"),
		d.boolProp("Visible", "Visible"),
		d.enumProp("Horizontal Anchor", "HorizontalAnchor",
			"Left", "Right", "Center", "Left+Right", "Left+Center", "Right+Center",
		),
		d.enumProp("Vertical Anchor", "VerticalAnchor",
			"Top", "Bottom", "Center", "Top+Bottom", "Top+Center", "Bottom+Center",
		),
		d.intProp("X", "X"),
		d.intProp("Y", "Y"),
		d.intProp("Width", "Width"),
		d.intProp("Height", "Height"),
		d.intProp("Inner X", "InnerX"),
		d.intProp("Inner Y", "InnerY"),
		d.intProp("Inner Width", "InnerWidth"),
		d.intProp("Inner Height", "InnerHeight"),
		d.enumProp("Alignment", "Alignment",
			"Left", "Center", "Right",
		),
		d.boolProp("Checked", "Checked"),
		d.intProp("Arrow Increment", "ArrowIncrement"),
		d.intProp("Mouse Increment", "MouseIncrement"),
		d.intProp("Min", "Min"),
		d.intProp("Max", "Max"),
		d.intProp("Value", "Value"),
		d.floatProp("Min", "Min"),
		d.floatProp("Max", "Max"),
		d.floatProp("Value", "Value"),
		d.intProp("Cursor Position", "CursorPosition"),
		d.intProp("Precision", "Precision", 1, 6),
		d.enumProp("Orientation", "Orientation",
			"Horizontal", "Vertical",
		),
		d.enumProp("Tick Position", "TickPosition",
			"Right/Bottom", "Left/Top", "Both Sides",
		),
		d.intProp("Tick Frequency", "TickFrequency"),
		d.boolProp("Ticks Visible", "TicksVisible"),
		d.enumProp("Border Style", "BorderStyle",
			"None", "Single Line", "Sunken", "Sunken Thick", "Raised",
		),
		d.intProp("Character Limit", "CharacterLimit", 1, 0x7FFFFFFE),
		d.boolProp("Is Password", "IsPassword"),
		d.boolProp("Read Only", "ReadOnly"),
		d.boolProp("Writes Tabs", "WritesTabs"),
		d.stringListProp("Items", "Items"),
		d.intProp("Selected Index", "SelectedIndex", -1, math.MaxInt32),
		d.boolProp("Vertical", "Vertical"),
		d.boolProp("Moves Forever", "MovesForever"),
		d.boolProp("Word Wrap", "WordWrap"),
	}

	fontProps := wui.NewPanel()
	w.Add(fontProps)
	useParentFont, useParentFontPanel := d.boolPanel(fontProps, "Use Parent Font")
	fontName, fontNamePanel := d.stringPanel(fontProps, "Name")
	fontName.SetCharacterLimit(31)
	fontHeight, fontHeightPanel := d.intPanel(fontProps, "Height")
	fontBold, fontBoldPanel := d.boolPanel(fontProps, "Bold")
	fontBold.SetFont(bold)
	fontItalic, fontItalicPanel := d.boolPanel(fontProps, "Italic")
	fontItalic.SetFont(italic)
	fontUnderlined, fontUnderlinedPanel := d.boolPanel(fontProps, "Underlined")
	fontUnderlined.SetFont(underlined)
	fontStrikedOut, fontStrikedOutPanel := d.boolPanel(fontProps, "StrikedOut")
	fontStrikedOut.SetFont(strikedOut)
	for _, p := range []*wui.Panel{
		useParentFontPanel,
		fontNamePanel,
		fontHeightPanel,
		fontBoldPanel,
		fontItalicPanel,
		fontUnderlinedPanel,
		fontStrikedOutPanel,
	} {
		p.SetX(p.X() - 40)
	}
	fontProps.SetBorderStyle(wui.PanelBorderSunken)
	{
		fontLabel := wui.NewLabel()
		fontLabel.SetText("Font")
		fontLabel.SetY(propMargin + 5)
		fontLabel.SetHeight(13)
		fontLabel.SetAlignment(wui.AlignCenter)
		fontLabel.SetFont(bold)
		fontProps.Add(fontLabel)

		y := fontLabel.Y() + fontLabel.Height() + 10
		for _, panel := range []*wui.Panel{
			useParentFontPanel,
			fontNamePanel,
			fontHeightPanel,
			fontBoldPanel,
			fontItalicPanel,
			fontUnderlinedPanel,
			fontStrikedOutPanel,
		} {
			panel.SetY(y)
			y += panel.Height()
		}
		fontProps.SetBounds(15, 0, 175, y+5)
		fontLabel.SetWidth(fontProps.InnerWidth())
	}
	updateFont := func() {
		f, ok := d.active.(fonter)
		if !ok {
			return
		}
		useParent := useParentFont.Checked()
		fontName.SetEnabled(!useParent)
		fontHeight.SetEnabled(!useParent)
		fontBold.SetEnabled(!useParent)
		fontItalic.SetEnabled(!useParent)
		fontUnderlined.SetEnabled(!useParent)
		fontStrikedOut.SetEnabled(!useParent)
		if useParent {
			f.SetFont(nil)
		} else {
			font, err := wui.NewFont(wui.FontDesc{
				Name:       fontName.Text(),
				Height:     fontHeight.Value(),
				Bold:       fontBold.Checked(),
				Italic:     fontItalic.Checked(),
				Underlined: fontUnderlined.Checked(),
				StrikedOut: fontStrikedOut.Checked(),
			})
			if err == nil {
				f.SetFont(font)
			}
		}
		d.preview.Paint()
	}
	useParentFont.SetOnChange(func(disable bool) { updateFont() })
	fontName.SetOnTextChange(func() { updateFont() })
	fontHeight.SetOnValueChange(func(int) { updateFont() })
	fontBold.SetOnChange(func(bool) { updateFont() })
	fontItalic.SetOnChange(func(bool) { updateFont() })
	fontUnderlined.SetOnChange(func(bool) { updateFont() })
	fontStrikedOut.SetOnChange(func(bool) { updateFont() })

	appIcon := w32.LoadIcon(0, w32.MakeIntResource(w32.IDI_APPLICATION))
	appIconWidth := w32.GetSystemMetrics(w32.SM_CXICON)
	appIconHeight := w32.GetSystemMetrics(w32.SM_CYICON)
	appIconWidth, appIconHeight = 17, 17

	defaultCursor := w.Cursor()

	leftSlider := wui.NewPanel()
	leftSlider.SetBounds(195, -1, 5, 602)
	leftSlider.SetBorderStyle(wui.PanelBorderSingleLine)
	leftSlider.SetVerticalAnchor(wui.AnchorMinAndMax)
	w.Add(leftSlider)

	rightSlider := wui.NewPanel()
	rightSlider.SetBounds(600, -1, 5, 602)
	rightSlider.SetBorderStyle(wui.PanelBorderSingleLine)
	rightSlider.SetVerticalAnchor(wui.AnchorMinAndMax)
	rightSlider.SetHorizontalAnchor(wui.AnchorMax)
	w.Add(rightSlider)

	panelTemplate := wui.NewPanel()
	panelTemplate.SetBounds(20, 10, 150, 50)
	panelTemplate.SetBorderStyle(wui.PanelBorderSingleLine)
	panelText := wui.NewLabel()
	panelText.SetText("Panel")
	panelText.SetAlignment(wui.AlignCenter)
	panelText.SetSize(panelTemplate.InnerWidth(), panelTemplate.InnerHeight())
	panelTemplate.Add(panelText)

	paintBoxTemplate := wui.NewPaintBox()
	paintBoxTemplate.SetBounds(20, 67, 150, 50)

	textEditTemplate := wui.NewTextEdit()
	textEditTemplate.SetBounds(20, 124, 150, 50)
	textEditTemplate.SetText("Text Edit")

	editLineTemplate := wui.NewEditLine()
	editLineTemplate.SetBounds(20, 181, 150, 20)
	editLineTemplate.SetText("Text Edit Line")

	comboTemplate := wui.NewComboBox()
	comboTemplate.SetBounds(20, 210, 150, 21)
	comboTemplate.AddItem("Combo Box")
	comboTemplate.SetSelectedIndex(0)

	sliderTemplate := wui.NewSlider()
	sliderTemplate.SetBounds(20, 245, 150, 45)

	progressTemplate := wui.NewProgressBar()
	progressTemplate.SetBounds(20, 295, 150, 25)
	progressTemplate.SetValue(0.5)

	buttonTemplate := wui.NewButton()
	buttonTemplate.SetText("Button")
	buttonTemplate.SetBounds(20, 329, 85, 25)

	intTemplate := wui.NewIntUpDown()
	intTemplate.SetBounds(20, 362, 80, 22)

	floatTemplate := wui.NewFloatUpDown()
	floatTemplate.SetBounds(20, 392, 80, 22)

	checkBoxTemplate := wui.NewCheckBox()
	checkBoxTemplate.SetText("Check Box")
	checkBoxTemplate.SetChecked(true)
	checkBoxTemplate.SetBounds(20, 423, 100, 17)

	radioButtonTemplate := wui.NewRadioButton()
	radioButtonTemplate.SetText("Radio Button")
	radioButtonTemplate.SetChecked(true)
	radioButtonTemplate.SetBounds(20, 448, 100, 17)

	labelTemplate := wui.NewLabel()
	labelTemplate.SetText("Text Label")
	labelTemplate.SetBounds(20, 473, 150, 13)

	allTemplates := []wui.Control{
		panelTemplate,
		paintBoxTemplate,
		textEditTemplate,
		editLineTemplate,
		comboTemplate,
		sliderTemplate,
		progressTemplate,
		buttonTemplate,
		intTemplate,
		floatTemplate,
		checkBoxTemplate,
		radioButtonTemplate,
		labelTemplate,
	}

	var highlightedTemplate, controlToAdd wui.Control
	var templateDx, templateDy int

	palette := wui.NewPaintBox()
	palette.SetBounds(605, 0, 195, 600)
	palette.SetHorizontalAnchor(wui.AnchorMax)
	palette.SetVerticalAnchor(wui.AnchorMinAndMax)
	palette.SetOnPaint(func(c *wui.Canvas) {
		w, h := c.Size()
		c.FillRect(0, 0, w, h, wui.RGB(240, 240, 240))
		for _, template := range allTemplates {
			drawControl(template, c)
		}
		// Highlight what is under the mouse.
		if highlightedTemplate != nil {
			x, y, w, h := highlightedTemplate.Bounds()
			c.DrawRect(x-1, y-1, w+2, h+2, wui.RGB(255, 0, 255))
			c.DrawRect(x-2, y-2, w+4, h+4, wui.RGB(255, 0, 255))
		}
	})
	palette.SetOnMouseMove(func(x, y int) {
		oldHighlight := highlightedTemplate
		highlightedTemplate = nil
		for _, c := range allTemplates {
			if contains(c, x, y) {
				highlightedTemplate = c
			}
		}
		if highlightedTemplate != oldHighlight {
			palette.Paint()
		}
	})
	w.Add(palette)

	nameText := wui.NewLabel()
	nameText.SetText("Variable Name")
	nameText.SetAlignment(wui.AlignRight)
	nameText.SetBounds(10, 10, 85, 20)
	w.Add(nameText)
	name := wui.NewEditLine()
	name.SetBounds(100, 10, 90, 22)
	w.Add(name)

	d.preview.SetBounds(200, 0, 400, 600)
	d.preview.SetHorizontalAnchor(wui.AnchorMinAndMax)
	d.preview.SetVerticalAnchor(wui.AnchorMinAndMax)
	white := wui.RGB(255, 255, 255)
	black := wui.RGB(0, 0, 0)

	editOnPaint := wui.NewButton()
	editOnPaint.SetText("OnPaint")
	editOnPaint.SetBounds(105, 500, 85, 25)
	editOnPaint.SetVisible(false) // TODO Bring this back.
	w.Add(editOnPaint)

	name.SetOnTextChange(func() {
		names[d.active] = name.Text()
	})
	editOnPaint.SetOnClick(func() {
		p, valid := d.active.(*wui.PaintBox)
		if !valid {
			panic("OnPaint only valid for paint boxes")
		}

		dlg := wui.NewWindow()
		dlg.SetPosition(w32.ClientToScreen(w32.HWND(d.preview.Handle()), 0, 0))
		dlg.SetSize(d.preview.Size())

		code := wui.NewTextEdit()
		font, _ := wui.NewFont(wui.FontDesc{Name: "Courier New", Height: -15})
		code.SetFont(font)
		code.SetWritesTabs(true)
		// TODO code.SetLineBreaks("\n")
		code.SetBounds(0, 0, dlg.InnerWidth(), dlg.InnerHeight()-30)
		code.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndMax)
		dlg.Add(code)

		onPaint := event{p, "OnPaint"}
		if events[onPaint] == "" {
			events[onPaint] = "func(canvas *wui.Canvas) {\n\t\n}"
		}
		code.SetText(strings.Replace(events[onPaint], "\n", "\r\n", -1))
		code.SetCursorPosition(29)

		ok := wui.NewButton()
		ok.SetText("OK")
		ok.SetBounds(dlg.InnerWidth()/2-87, dlg.InnerHeight()-28, 85, 25)
		ok.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
		ok.SetOnClick(func() {
			events[onPaint] = strings.Replace(code.Text(), "\r", "", -1)
			dlg.Close()
		})
		dlg.Add(ok)

		cancel := wui.NewButton()
		cancel.SetText("Cancel")
		cancel.SetBounds(dlg.InnerWidth()/2+2, dlg.InnerHeight()-28, 85, 25)
		cancel.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
		cancel.SetOnClick(dlg.Close)
		dlg.Add(cancel)

		dlg.SetOnShow(code.Focus)

		dlg.ShowModal()
	})

	// updateProperties refreshes the visible UI properties by reading the
	// values in from the active control.
	d.updateProperties = func() {
		for _, prop := range uiProps {
			if prop.panel.Visible() {
				prop.update()
			}
		}
	}

	activate := func(newActive node) {
		d.active = newActive

		name.SetText(names[d.active])
		y := name.Y() + name.Height() + propMargin

		for _, prop := range uiProps {
			m, hasProp := reflect.TypeOf(d.active).MethodByName(prop.setter)
			show := hasProp && prop.rightType(m.Type.In(1))
			prop.panel.SetVisible(show)
			if show {
				prop.panel.SetY(y)
				y += prop.panel.Height()
			}
		}
		d.updateProperties()

		f, hasFont := d.active.(fonter)
		fontProps.SetVisible(hasFont)
		if hasFont {
			fontProps.SetY(y)
			y += fontProps.Height()
			font := f.Font()
			if _, isWindow := d.active.(*wui.Window); isWindow {
				useParentFont.SetEnabled(false)
				useParentFont.SetChecked(false)
			} else {
				useParentFont.SetEnabled(true)
				useParentFont.SetChecked(font == nil)
			}
			if font != nil {
				fontName.SetText(font.Desc.Name)
				fontHeight.SetValue(font.Desc.Height)
				fontBold.SetChecked(font.Desc.Bold)
				fontItalic.SetChecked(font.Desc.Italic)
				fontUnderlined.SetChecked(font.Desc.Underlined)
				fontStrikedOut.SetChecked(font.Desc.StrikedOut)
			}
		}
	}
	activate(theWindow)

	// mouseMode constants.
	const (
		idleMouse = iota
		addControl
		dragTopLeft
		dragTop
		dragTopRight
		dragRight
		dragBottomRight
		dragBottom
		dragBottomLeft
		dragLeft
		dragAll
	)
	var (
		mouseMode         = idleMouse
		nextDragMouseMode int
		nextToDrag        node
	)
	dragging := func() bool {
		return dragTopLeft <= mouseMode && mouseMode <= dragAll
	}

	var xOffset, yOffset int
	d.preview.SetOnPaint(func(c *wui.Canvas) {
		// Place the inner top-left at 20,40.
		xOffset = 20 - (theWindow.InnerX() - theWindow.X())
		yOffset = 40 - (theWindow.InnerY() - theWindow.Y())
		width, height := theWindow.Size()
		innerWidth, innerHeight := theWindow.InnerSize()
		borderSize := (width - innerWidth) / 2
		topBorderSize := height - borderSize - innerHeight
		innerX = xOffset + borderSize
		innerY = yOffset + topBorderSize
		inner := makeOffsetDrawer(c, innerX, innerY)

		c.FillRect(0, 0, d.preview.Width(), d.preview.Height(), white)

		// Clear inner area.
		c.FillRect(innerX, innerY, innerWidth, innerHeight, wui.RGB(240, 240, 240))

		// Draw all the window contents.
		drawContainer(theWindow, inner)

		// Draw the window border, icon and title.
		borderColor := wui.RGB(100, 200, 255)
		c.FillRect(xOffset, yOffset, width, topBorderSize, borderColor)
		c.FillRect(xOffset, yOffset, borderSize, height, borderColor)
		c.FillRect(xOffset, yOffset+height-borderSize, width, borderSize, borderColor)
		c.FillRect(xOffset+width-borderSize, yOffset, borderSize, height, borderColor)

		if theWindow.HasBorder() {
			_, textH := c.TextExtent(theWindow.Title())
			c.TextOut(
				xOffset+borderSize+appIconWidth+5,
				yOffset+(topBorderSize-textH)/2,
				theWindow.Title(),
				black,
			)

			w := topBorderSize
			h := w - 8
			y := yOffset + 4
			right := xOffset + width - borderSize
			x0 := right - 3*w - 2
			x1 := right - 2*w - 1
			x2 := right - 1*w - 0
			iconSize := h / 2
			if theWindow.HasMinButton() || theWindow.HasMaxButton() {
				{
					// Minimize button.
					c.FillRect(x0, y, w, h, wui.RGB(240, 240, 240))
					cx := x0 + (w-iconSize)/2
					cy := y + h - 1 - (iconSize+1)/2
					color := black
					if !theWindow.HasMinButton() {
						color = wui.RGB(204, 204, 204)
					}
					c.Line(cx, cy, cx+iconSize, cy, color)
				}
				{
					// Maximize button.
					c.FillRect(x1, y, w, h, wui.RGB(240, 240, 240))
					cx := x1 + (w-iconSize)/2
					cy := y + (h-iconSize)/2
					color := black
					if !theWindow.HasMaxButton() {
						color = wui.RGB(204, 204, 204)
					}
					c.DrawRect(cx, cy, iconSize, iconSize, color)
				}
			}
			// Close button.
			color := black
			backColor := wui.RGB(255, 128, 128)
			if !theWindow.HasCloseButton() {
				color = wui.RGB(204, 204, 204)
				backColor = wui.RGB(240, 240, 240)
			}
			c.FillRect(x2, y, w, h, backColor)
			cx := x2 + (w-iconSize)/2
			cy := y + (h-iconSize)/2
			c.Line(cx, cy, cx+iconSize, cy+iconSize, color)
			c.Line(cx, cy+iconSize-1, cx+iconSize, cy-1, color)

			w32.DrawIconEx(
				w32.HDC(c.Handle()),
				xOffset+borderSize,
				yOffset+(topBorderSize-appIconHeight)/2,
				appIcon,
				appIconWidth, appIconHeight,
				0, 0, w32.DI_NORMAL,
			)
		}

		// Highlight the currently selected child control, except if dragging
		// it with the mouse.
		if !dragging() && d.active != nil && d.active != theWindow {
			x, y, w, h := d.active.Bounds()
			parent := d.active.Parent()
			for parent != theWindow {
				dx, dy, _, _ := parent.InnerBounds()
				x += dx
				y += dy
				parent = parent.Parent()
			}
			w = max(w, 0)
			h = max(h, 0)
			inner.DrawRect(x, y, w, h, wui.RGB(255, 0, 255))
			inner.DrawRect(x+1, y+1, w-2, h-2, wui.RGB(255, 0, 255))
		}

		if controlToAdd != nil {
			drawControl(controlToAdd, c)
		}
	})
	w.Add(d.preview)

	var (
		dragStartX, dragStartY                                  int
		preResizeX, preResizeY, preResizeWidth, preResizeHeight int
	)

	lastX, lastY := -999999, -999999
	w.SetOnMouseMove(func(x, y int) {
		if x == lastX && y == lastY {
			return
		}
		lastX, lastY = x, y

		if mouseMode == addControl {
			if contains(d.preview, x, y) {
				_, _, w, h := controlToAdd.Bounds()
				relX := x - d.preview.X()
				relY := y - d.preview.Y()
				if false {
					// TODO Align to some nice-looking grid unless Ctrl is held
					// down for example. NOTE that this right now contains a
					// bug, relX is not in window client coordinates, it is
					// relative to the preview paint box and thus we can never
					// get to 0,0 with this.
					const gridSize = 10
					relX = relX / gridSize * gridSize
					relY = relY / gridSize * gridSize
				}
				relX += templateDx
				relY += templateDy
				controlToAdd.SetBounds(relX, relY, w, h)
			}
			d.preview.Paint()
		} else if mouseMode == idleMouse {
			// See if the cursor is over the edge of the active control. In that
			// case show the resize cursor and remember what to resize and in
			// which direction.
			x -= d.preview.X()
			y -= d.preview.Y()
			x -= xOffset
			y -= yOffset
			nextToDrag = d.active
			ax, ay, aw, ah := relativeBounds(d.active, theWindow)
			const margin = 6
			corner := func(x, y int) rectangle {
				return rect(x-margin, y-margin, 2*margin, 2*margin)
			}
			var (
				// These are the draggable areas for the active control.
				topLeft     = corner(ax, ay)
				top         = rect(ax, ay-margin, aw, 2*margin)
				topRight    = corner(ax+aw, ay)
				right       = rect(ax+aw-margin, ay, 2*margin, ah)
				bottomRight = corner(ax+aw, ay+ah)
				bottom      = rect(ax, ay+ah-margin, aw, 2*margin)
				bottomLeft  = corner(ax, ay+ah)
				left        = rect(ax-margin, ay, 2*margin, ah)
			)
			if d.active == theWindow {
				// The main window can only be dragged right and bottom so we
				// reset the other drag areas. They will not be triggered for
				// the main window.
				topLeft = rectangle{}
				top = rectangle{}
				topRight = rectangle{}
				bottomLeft = rectangle{}
				left = rectangle{}
			}
			var (
				// No matter what the active control is, we want to be able to
				// drag the main window always so we check for that separately.
				winX, winY, winW, winH = relativeBounds(theWindow, theWindow)
				winRight               = rect(winX+winW-margin, winY, 2*margin, winH)
				winBottomRight         = corner(winX+winW, winY+winH)
				winBottom              = rect(winX, winY+winH-margin, winW, 2*margin)
			)
			if winBottomRight.contains(x, y) {
				nextDragMouseMode = dragBottomRight
				w.SetCursor(wui.CursorSizeNWSE)
				nextToDrag = theWindow
			} else if winRight.contains(x, y) {
				nextDragMouseMode = dragRight
				w.SetCursor(wui.CursorSizeWE)
				nextToDrag = theWindow
			} else if winBottom.contains(x, y) {
				nextDragMouseMode = dragBottom
				w.SetCursor(wui.CursorSizeNS)
				nextToDrag = theWindow
			} else if topLeft.contains(x, y) {
				nextDragMouseMode = dragTopLeft
				w.SetCursor(wui.CursorSizeNWSE)
			} else if topRight.contains(x, y) {
				nextDragMouseMode = dragTopRight
				w.SetCursor(wui.CursorSizeNESW)
			} else if bottomRight.contains(x, y) {
				nextDragMouseMode = dragBottomRight
				w.SetCursor(wui.CursorSizeNWSE)
			} else if bottomLeft.contains(x, y) {
				nextDragMouseMode = dragBottomLeft
				w.SetCursor(wui.CursorSizeNESW)
			} else if top.contains(x, y) {
				nextDragMouseMode = dragTop
				w.SetCursor(wui.CursorSizeNS)
			} else if right.contains(x, y) {
				nextDragMouseMode = dragRight
				w.SetCursor(wui.CursorSizeWE)
			} else if bottom.contains(x, y) {
				nextDragMouseMode = dragBottom
				w.SetCursor(wui.CursorSizeNS)
			} else if left.contains(x, y) {
				nextDragMouseMode = dragLeft
				w.SetCursor(wui.CursorSizeWE)
			} else {
				// If we are not over a draggable border but inside the active
				// control, we drag it completely without resizing.
				// If we are not inside the active control we activate another
				// control with this mouse click.
				// There is a catch, though. If the active control is a
				// container and we are over a child control, we do not want to
				// drag the container, we want to activate the child control,
				// even though the mouse is still inside the active container.
				// Then there is an exception to this which is the main window.
				// We cannot drag that as a whole at all.
				innerX, innerY, _, _ := theWindow.InnerBounds()
				outerX, outerY, _, _ := theWindow.Bounds()
				relX := x - (innerX - outerX)
				relY := y - (innerY - outerY)
				if theWindow != d.active &&
					d.active == findControlAt(theWindow, relX, relY) {
					nextDragMouseMode = dragAll
					w.SetCursor(wui.CursorSizeAll)
				} else {
					nextDragMouseMode = idleMouse
					w.SetCursor(defaultCursor)
				}
			}
		} else {
			// In this case we are dragging, update the relevant parts of the
			// control being dragged.
			dx := x - dragStartX
			dy := y - dragStartY
			x, y, w, h := preResizeX, preResizeY, preResizeWidth, preResizeHeight
			switch mouseMode {
			case dragTopLeft:
				dx = min(dx, w)
				dy = min(dy, h)
				nextToDrag.SetBounds(x+dx, y+dy, w-dx, h-dy)
			case dragTop:
				dy = min(dy, h)
				nextToDrag.SetBounds(x, y+dy, w, h-dy)
			case dragTopRight:
				dx = max(dx, -w)
				dy = min(dy, h)
				nextToDrag.SetBounds(x, y+dy, w+dx, h-dy)
			case dragRight:
				dx = max(dx, -w)
				nextToDrag.SetBounds(x, y, w+dx, h)
			case dragBottomRight:
				dx = max(dx, -w)
				dy = max(dy, -h)
				nextToDrag.SetBounds(x, y, w+dx, h+dy)
			case dragBottom:
				dy = max(dy, -h)
				nextToDrag.SetBounds(x, y, w, h+dy)
			case dragBottomLeft:
				dx = min(dx, w)
				dy = max(dy, -h)
				nextToDrag.SetBounds(x+dx, y, w-dx, h+dy)
			case dragLeft:
				dx = min(dx, w)
				nextToDrag.SetBounds(x+dx, y, w-dx, h)
			case dragAll:
				nextToDrag.SetBounds(x+dx, y+dy, w, h)
			}
			d.updateProperties()
			d.preview.Paint()
		}
	})

	w.SetOnMouseDown(func(button wui.MouseButton, x, y int) {
		if button == wui.MouseButtonLeft {
			if contains(palette, x, y) && highlightedTemplate != nil {
				controlToAdd = cloneControl(highlightedTemplate)
				hx, hy, _, _ := highlightedTemplate.Bounds()
				templateDx = hx - (x - palette.X())
				templateDy = hy - (y - palette.Y())
				mouseMode = addControl
				activate(theWindow)
				d.preview.Paint()
			} else if mouseMode == addControl {
				innerX, innerY, _, _ := theWindow.InnerBounds()
				outerX, outerY, _, _ := theWindow.Bounds()
				x, y, w, h := controlToAdd.Bounds()
				relX := x - (xOffset + innerX - outerX)
				relY := y - (yOffset + innerY - outerY)
				// Find the sub-container that this is to be placed in. Use the
				// center of the new control to determine where to add it.
				addToThis, x, y := findContainerAt(theWindow, relX+w/2, relY+h/2)
				controlToAdd.SetBounds(x-w/2, y-h/2, w, h)
				names[controlToAdd] = defaultName(controlToAdd)
				addToThis.Add(controlToAdd)
				activate(controlToAdd)
				controlToAdd = nil
				mouseMode = idleMouse
				name.Focus()
				name.SelectAll()
				d.preview.Paint()
			} else {
				dragStartX = x
				dragStartY = y
				preResizeX, preResizeY, preResizeWidth, preResizeHeight = nextToDrag.Bounds()
				mouseMode = nextDragMouseMode
				if mouseMode == idleMouse && contains(d.preview, x, y) {
					newActive := findControlAt(
						theWindow,
						x-d.preview.X()-innerX,
						y-d.preview.Y()-innerY,
					)
					if newActive != d.active {
						activate(newActive)
					}
				}
				d.preview.Paint()
			}
		}
	})

	w.SetOnMouseUp(func(button wui.MouseButton, x, y int) {
		if button == wui.MouseButtonLeft {
			if mouseMode != addControl {
				mouseMode = idleMouse
			}
		}
		// TODO Why does this not work?
		//nextDragMouseMode = idleMouse
		//w.OnMouseMove()(x+1, y)
		//w.OnMouseMove()(x, y)
		d.preview.Paint()
	})

	workingPath := ""
	setWorkingPath := func(path string) {
		workingPath = path
		title := "wui Designer"
		if path != "" {
			title += " - " + path
		}
		w.SetTitle(title)
	}
	setWorkingPath("")

	fileOpenMenu.SetOnClick(func() {
		open := wui.NewFileOpenDialog()
		open.SetTitle("Select a Go file containing one or more wui.Windows")
		open.AddFilter("Go file", ".go")
		if accept, path := open.ExecuteSingleSelection(w); accept {
			setWorkingPath(path)
			//wui.MessageBoxError("TODO", "Open is not yet implemented")
			if err := parseFile(theWindow, path); err != nil {
				wui.MessageBoxError("Error", err.Error())
			} else {
				d.preview.Paint()
			}
		}
	})

	saveCodeTo := func(path string) {
		code := generateCode(theWindow, false)
		err := ioutil.WriteFile(path, code, 0666)
		if err != nil {
			wui.MessageBoxError("Error", err.Error())
		} else {
			workingPath = path
		}
	}

	fileSaveAsMenu.SetOnClick(func() {
		save := wui.NewFileSaveDialog()
		save.SetAppendExt(true)
		save.AddFilter("Go file", ".go")
		if accept, path := save.Execute(w); accept {
			saveCodeTo(path)
			w32.ShellExecute(0, "open", path, "", "", w32.SW_SHOWNORMAL)
		}
	})

	fileSaveMenu.SetOnClick(func() {
		if workingPath != "" {
			saveCodeTo(workingPath)
		} else {
			fileSaveAsMenu.OnClick()()
		}
	})

	previewMenu.SetOnClick(func() {
		// We place the window such that it lies exactly over our drawing.
		x, y := w32.ClientToScreen(w32.HWND(w.Handle()), d.preview.X(), d.preview.Y())
		showPreview(w, theWindow, x+xOffset, y+yOffset)
	})

	exitMenu.SetOnClick(w.Close)

	// TODO Build undo/redo.
	//undoMenu.SetOnClick(func() {})
	//redoMenu.SetOnClick(func() {})

	deleteMenu.SetOnClick(func() {
		if d.active != nil && d.active != theWindow {
			c := d.active.(wui.Control)
			p := d.active.Parent()
			activate(p)
			p.Remove(c)
			d.preview.Paint()
		}
	})

	w.SetShortcut(fileOpenMenu.OnClick(), wui.KeyControl, wui.KeyO)
	w.SetShortcut(fileSaveMenu.OnClick(), wui.KeyControl, wui.KeyS)
	w.SetShortcut(fileSaveAsMenu.OnClick(), wui.KeyControl, wui.KeyShift, wui.KeyS)
	w.SetShortcut(previewMenu.OnClick(), wui.KeyControl, wui.KeyR)
	w.SetShortcut(undoMenu.OnClick(), wui.KeyControl, wui.KeyZ)
	w.SetShortcut(redoMenu.OnClick(), wui.KeyControl, wui.KeyShift, wui.KeyZ)
	w.SetShortcut(deleteMenu.OnClick(), wui.KeyControl, wui.KeyDelete)

	//w.SetShortcut(w.Close, wui.KeyEscape) // TODO ESC for debugging

	w.SetState(wui.WindowMaximized)
	w.Show()
}

func rect(x, y, width, height int) rectangle {
	return rectangle{x: x, y: y, w: width, h: height}
}

type rectangle struct {
	x, y, w, h int
}

func (r rectangle) contains(x, y int) bool {
	return x >= r.x && y >= r.y && x < r.x+r.w && y < r.y+r.h
}

func defaultWindow() *wui.Window {
	font, _ := wui.NewFont(wui.FontDesc{Name: "Tahoma", Height: -11})
	w := wui.NewWindow()
	w.SetFont(font)
	w.SetTitle("Window")
	return w
}

func findControlAt(parent wui.Container, x, y int) node {
	for _, child := range parent.Children() {
		if contains(child, x, y) {
			if container, ok := child.(wui.Container); ok {
				dx, dy, _, _ := container.Bounds()
				return findControlAt(container, x-dx, y-dy)
			}
			return child
		}
	}
	return parent
}

func contains(b bounder, atX, atY int) bool {
	x, y, w, h := b.Bounds()
	return atX >= x && atY >= y && atX < x+w && atY < y+h
}

type bounder interface {
	Bounds() (x, y, width, height int)
}

func innerContains(b innerBounder, atX, atY int) bool {
	x, y, w, h := b.InnerBounds()
	return atX >= x && atY >= y && atX < x+w && atY < y+h
}

type innerBounder interface {
	InnerBounds() (x, y, width, height int)
}

type drawer interface {
	PushDrawRegion(x, y, width, height int)
	PopDrawRegion()
	Line(x1, y1, x2, y2 int, color wui.Color)
	DrawRect(x, y, w, h int, color wui.Color)
	FillRect(x, y, w, h int, color wui.Color)
	DrawEllipse(x, y, w, h int, color wui.Color)
	FillEllipse(x, y, w, h int, color wui.Color)
	TextRectFormat(x, y, w, h int, s string, format wui.Format, color wui.Color)
	TextExtent(s string) (width, height int)
	TextOut(x, y int, s string, color wui.Color)
	Polygon(p []wui.Point, color wui.Color)
	SetFont(*wui.Font)
}

func makeOffsetDrawer(base drawer, dx, dy int) drawer {
	return &offsetDrawer{base: base, dx: dx, dy: dy}
}

type offsetDrawer struct {
	base   drawer
	dx, dy int
}

func (d *offsetDrawer) PushDrawRegion(x, y, width, height int) {
	d.base.PushDrawRegion(x+d.dx, y+d.dy, width, height)
}

func (d *offsetDrawer) PopDrawRegion() {
	d.base.PopDrawRegion()
}

func (d *offsetDrawer) DrawRect(x, y, w, h int, color wui.Color) {
	d.base.DrawRect(x+d.dx, y+d.dy, w, h, color)
}

func (d *offsetDrawer) FillRect(x, y, w, h int, color wui.Color) {
	d.base.FillRect(x+d.dx, y+d.dy, w, h, color)
}

func (d *offsetDrawer) DrawEllipse(x, y, w, h int, color wui.Color) {
	d.base.DrawEllipse(x+d.dx, y+d.dy, w, h, color)
}

func (d *offsetDrawer) FillEllipse(x, y, w, h int, color wui.Color) {
	d.base.FillEllipse(x+d.dx, y+d.dy, w, h, color)
}

func (d *offsetDrawer) TextRectFormat(
	x, y, w, h int, s string, format wui.Format, color wui.Color,
) {
	d.base.TextRectFormat(x+d.dx, y+d.dy, w, h, s, format, color)
}

func (d *offsetDrawer) TextExtent(s string) (width, height int) {
	return d.base.TextExtent(s)
}

func (d *offsetDrawer) TextOut(x, y int, s string, color wui.Color) {
	d.base.TextOut(x+d.dx, y+d.dy, s, color)
}

func (d *offsetDrawer) Polygon(p []wui.Point, color wui.Color) {
	for i := range p {
		p[i].X += int32(d.dx)
		p[i].Y += int32(d.dy)
	}
	d.base.Polygon(p, color)
}

func (d *offsetDrawer) Line(x1, y1, x2, y2 int, color wui.Color) {
	d.base.Line(x1+d.dx, y1+d.dy, x2+d.dx, y2+d.dy, color)
}

func (d *offsetDrawer) SetFont(f *wui.Font) {
	d.base.SetFont(f)
}

type node interface {
	Parent() wui.Container
	Bounds() (x, y, width, height int)
	SetBounds(x, y, width, height int)
}

func showPreview(parent, w *wui.Window, x, y int) {
	// Create a centered progress dialog that cannot be closed until the preview
	// is shown.
	canClose := make(chan bool, 1)
	progress := wui.NewWindow()
	progress.SetHasMinButton(false)
	progress.SetHasMaxButton(false)
	progress.SetHasCloseButton(false)
	progress.SetResizable(false)
	progress.SetTitle("Generating Preview...")
	progress.DisableAltF4()
	progress.SetOnCanClose(func() bool {
		return <-canClose
	})
	progress.SetInnerSize(420, 50)
	progress.SetX(parent.X() + (parent.Width()-progress.Width())/2)
	progress.SetY(parent.Y() + (parent.Height()-progress.Height())/2)
	p := wui.NewProgressBar()
	p.SetMovesForever(true)
	p.SetBounds(10, 10, 400, 30)
	progress.Add(p)

	// Generate the code in a different go routine while the progress bar is
	// showing.
	go func() {
		defer func() {
			canClose <- true
			progress.Close()
		}()

		// For the preview we set a temporary window position to align it with
		// the preview shown in the designer.
		oldX, oldY := w.Position()
		w.SetPosition(x, y)
		code := generateCode(w, true)
		w.SetPosition(oldX, oldY)

		// Write the Go file to our temporary build dir.
		goFile := filepath.Join(buildDir, "wui_designer_temp_file.go")
		err := ioutil.WriteFile(goFile, code, 0666)
		if err != nil {
			wui.MessageBoxError("Error", err.Error())
			return
		}
		defer os.Remove(goFile)

		// Build the executable into our temporary build dir.
		exeFile := filepath.Join(buildDir, buildPrefix+strconv.Itoa(buildCount)+".exe")
		buildCount++

		// Do the build synchronously and report any build errors.
		var tryOutput []byte
		var tryErr error
		try := func(cmd string, args ...string) {
			if err != nil {
				return
			}
			command := exec.Command(cmd, args...)
			command.Dir = buildDir
			tryOutput, tryErr = command.CombinedOutput()
		}
		try("go", "mod", "init", "temp/wui/preview")
		try("go", "mod", "tidy")
		try("go", "build", "-o", exeFile, goFile)
		if tryErr != nil {
			wui.MessageBoxError("Error", tryErr.Error()+"\r\n"+string(tryOutput))
			return
		}

		// Start the program in parallel so we can have multiple previews open at
		// once.
		exec.Command(exeFile).Start()
	}()

	progress.ShowModal()
}

func cloneControl(c wui.Control) wui.Control {
	// TODO Use the properties for this.
	switch x := c.(type) {
	case *wui.Button:
		b := wui.NewButton()
		b.SetText(x.Text())
		b.SetBounds(0, 0, x.Width(), x.Height())
		return b
	case *wui.CheckBox:
		c := wui.NewCheckBox()
		c.SetText(x.Text())
		c.SetChecked(x.Checked())
		c.SetBounds(0, 0, x.Width(), x.Height())
		return c
	case *wui.RadioButton:
		r := wui.NewRadioButton()
		r.SetText(x.Text())
		r.SetBounds(0, 0, x.Width(), x.Height())
		return r
	case *wui.Slider:
		s := wui.NewSlider()
		s.SetMinMax(x.MinMax())
		s.SetCursorPosition(x.CursorPosition())
		s.SetTickFrequency(x.TickFrequency())
		s.SetArrowIncrement(x.ArrowIncrement())
		s.SetMouseIncrement(x.MouseIncrement())
		s.SetTicksVisible(x.TicksVisible())
		s.SetOrientation(x.Orientation())
		s.SetTickPosition(x.TickPosition())
		s.SetBounds(0, 0, x.Width(), x.Height())
		return s
	case *wui.Panel:
		p := wui.NewPanel()
		p.SetBorderStyle(x.BorderStyle())
		p.SetBounds(0, 0, x.Width(), x.Height())
		return p
	case *wui.Label:
		l := wui.NewLabel()
		l.SetText(x.Text())
		l.SetAlignment(x.Alignment())
		l.SetBounds(0, 0, x.Width(), x.Height())
		return l
	case *wui.PaintBox:
		p := wui.NewPaintBox()
		p.SetBounds(0, 0, x.Width(), x.Height())
		return p
	case *wui.EditLine:
		e := wui.NewEditLine()
		e.SetBounds(0, 0, x.Width(), x.Height())
		e.SetText(x.Text())
		e.SetIsPassword(x.IsPassword())
		e.SetCharacterLimit(x.CharacterLimit())
		e.SetReadOnly(x.ReadOnly())
		return e
	case *wui.IntUpDown:
		n := wui.NewIntUpDown()
		n.SetBounds(0, 0, x.Width(), x.Height())
		n.SetMinMax(x.MinMax())
		n.SetValue(x.Value())
		return n
	case *wui.ComboBox:
		c := wui.NewComboBox()
		c.SetItems(x.Items())
		c.SetSelectedIndex(x.SelectedIndex())
		c.SetBounds(0, 0, x.Width(), x.Height())
		return c
	case *wui.ProgressBar:
		p := wui.NewProgressBar()
		p.SetValue(x.Value())
		p.SetVertical(x.Vertical())
		p.SetMovesForever(x.MovesForever())
		p.SetBounds(0, 0, x.Width(), x.Height())
		return p
	case *wui.FloatUpDown:
		f := wui.NewFloatUpDown()
		f.SetBounds(0, 0, x.Width(), x.Height())
		f.SetMinMax(x.MinMax())
		f.SetPrecision(x.Precision())
		f.SetValue(x.Value())
		return f
	case *wui.TextEdit:
		t := wui.NewTextEdit()
		t.SetBounds(0, 0, x.Width(), x.Height())
		t.SetCharacterLimit(x.CharacterLimit())
		t.SetWordWrap(x.WordWrap())
		t.SetText(x.Text())
		return t
	default:
		panic("unhandled control type in cloneControl")
	}
}

type fonter interface {
	Font() *wui.Font
	SetFont(*wui.Font)
}

func findContainerAt(c wui.Container, x, y int) (innerMost wui.Container, atX, atY int) {
	for _, child := range c.Children() {
		if container, ok := child.(wui.Container); ok {
			if innerContains(container, x, y) {
				dx, dy, _, _ := container.InnerBounds()
				return findContainerAt(container, x-dx, y-dy)
			}
		}
	}
	return c, x, y
}

// removeEmptyStrings changes the given input slice.
func removeEmptyStrings(items []string) []string {
	n := 0
	for i := range items {
		if items[i] != "" {
			items[n] = items[i]
			n++
		}
	}
	return items[:n]
}

func defaultName(of interface{}) string {
	typ := strings.TrimPrefix(reflect.TypeOf(of).String(), "*wui.")
	prefix := decapitalize(typ)
	i := 1
	for {
		name := prefix + strconv.Itoa(i)
		if !nameUsed(name) {
			return name
		}
		i++
	}
}

func decapitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

func nameUsed(name string) bool {
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}

type fontControl interface {
	Font() *wui.Font
	Parent() wui.Container
}

func getFont(f fontControl) *wui.Font {
	if f == nil {
		return nil
	}
	font := f.Font()
	if font != nil {
		return font
	}
	return getFont(f.Parent())
}

// relativeBounds returns a control's outer bounds relative to the outer bounds
// of the given container. If the control is the same as the container this will
// result in x and y being 0.
func relativeBounds(of node, in wui.Container) (x, y, width, height int) {
	x, y, width, height = of.Bounds()
	parent := of.Parent()
	for parent != nil {
		innerX, innerY, _, _ := parent.InnerBounds()
		x += innerX
		y += innerY
		parent = parent.Parent()
	}
	dx, dy, _, _ := in.Bounds()
	x -= dx
	y -= dy
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
