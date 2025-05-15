package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gonutz/wui/v2"
)

func drawContainer(container wui.Container, d drawer) {
	_, _, w, h := container.InnerBounds()
	d.PushDrawRegion(0, 0, w, h)
	for _, child := range container.Children() {
		if f, ok := child.(fontControl); ok {
			d.SetFont(getFont(f))
		}
		drawControl(child, d)
	}
	d.PopDrawRegion()
}

func drawControl(c wui.Control, d drawer) {
	switch x := c.(type) {
	case *wui.Button:
		drawButton(x, d)
	case *wui.RadioButton:
		drawRadioButton(x, d)
	case *wui.CheckBox:
		drawCheckBox(x, d)
	case *wui.Panel:
		drawPanel(x, d)
	case *wui.Slider:
		drawSlider(x, d)
	case *wui.Label:
		drawLabel(x, d)
	case *wui.PaintBox:
		drawPaintBox(x, d)
	case *wui.EditLine:
		drawEditLine(x, d)
	case *wui.IntUpDown:
		drawIntUpDown(x, d)
	case *wui.ComboBox:
		drawComboBox(x, d)
	case *wui.ProgressBar:
		drawProgressBar(x, d)
	case *wui.FloatUpDown:
		drawFloatUpDown(x, d)
	case *wui.TextEdit:
		drawTextEdit(x, d)
	default:
		panic("unhandled control type")
	}
}

func drawButton(b *wui.Button, d drawer) {
	x, y, w, h := b.Bounds()
	if w > 0 && h > 0 {
		d.DrawRect(x, y, w, h, wui.RGB(240, 240, 240))
	}
	if w > 2 && h > 2 {
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(173, 173, 173))
	}
	if w > 4 && h > 4 {
		d.FillRect(x+2, y+2, w-4, h-4, wui.RGB(225, 225, 225))
	}
	if w > 6 && h > 6 {
		d.SetFont(getFont(b))
		textW, textH := d.TextExtent(b.Text())
		d.PushDrawRegion(x+3, y+3, w-6, h-6)
		d.TextOut(x+(w-textW)/2, y+(h-textH)/2, b.Text(), wui.RGB(0, 0, 0))
		d.PopDrawRegion()
	}
}

func drawRadioButton(r *wui.RadioButton, d drawer) {
	x, y, w, h := r.Bounds()
	d.PushDrawRegion(x, y, w, h)
	d.FillRect(x, y, w, h, wui.RGB(240, 240, 240))
	d.FillEllipse(x, y+(h-13)/2, 13, 13, wui.RGB(255, 255, 255))
	d.DrawEllipse(x, y+(h-13)/2, 13, 13, wui.RGB(0, 0, 0))
	if r.Checked() {
		d.FillEllipse(x+3, y+(h-13)/2+3, 7, 7, wui.RGB(0, 0, 0))
	}
	_, textH := d.TextExtent(r.Text())
	d.TextOut(x+16, y+(h-textH)/2, r.Text(), wui.RGB(0, 0, 0))
	d.PopDrawRegion()
}

func drawCheckBox(c *wui.CheckBox, d drawer) {
	x, y, w, h := c.Bounds()
	d.PushDrawRegion(x, y, w, h)
	d.FillRect(x, y, w, h, wui.RGB(240, 240, 240))
	d.FillRect(x, y+(h-13)/2, 13, 13, wui.RGB(255, 255, 255))
	d.DrawRect(x, y+(h-13)/2, 13, 13, wui.RGB(0, 0, 0))
	if c.Checked() {
		// Draw two lines for the check mark. ✓
		startX := x + 2
		startY := y + (h-13)/2 + 6
		d.Line(startX, startY, startX+3, startY+3, wui.RGB(0, 0, 0))
		d.Line(startX+3, startY+2, startX+9, startY-4, wui.RGB(0, 0, 0))
	}
	_, textH := d.TextExtent(c.Text())
	d.TextOut(x+16, y+(h-textH)/2, c.Text(), wui.RGB(0, 0, 0))
	d.PopDrawRegion()
}

func drawPanel(p *wui.Panel, d drawer) {
	x, y, w, h := p.Bounds()
	if w <= 0 || h <= 0 {
		return
	}
	switch p.BorderStyle() {
	case wui.PanelBorderNone:
		d.DrawRect(x, y, w, h, wui.RGB(230, 230, 230))
	case wui.PanelBorderSingleLine:
		d.DrawRect(x, y, w, h, wui.RGB(100, 100, 100))
	case wui.PanelBorderRaised:
		d.Line(x, y, x+w, y, wui.RGB(227, 227, 227))
		d.Line(x, y, x, y+h, wui.RGB(227, 227, 227))
		d.Line(x+w-1, y, x+w-1, y+h, wui.RGB(105, 105, 105))
		d.Line(x, y+h-1, x+w, y+h-1, wui.RGB(105, 105, 105))
		d.Line(x+1, y+1, x+w-1, y+1, wui.RGB(255, 255, 255))
		d.Line(x+1, y+1, x+1, y+h-1, wui.RGB(255, 255, 255))
		d.Line(x+w-2, y+1, x+w-2, y+h-1, wui.RGB(160, 160, 160))
		d.Line(x+1, y+h-2, x+w-1, y+h-2, wui.RGB(160, 160, 160))
	case wui.PanelBorderSunken:
		d.Line(x, y, x+w, y, wui.RGB(160, 160, 160))
		d.Line(x, y, x, y+h, wui.RGB(160, 160, 160))
		d.Line(x+w-1, y, x+w-1, y+h, wui.RGB(255, 255, 255))
		d.Line(x, y+h-1, x+w, y+h-1, wui.RGB(255, 255, 255))
	case wui.PanelBorderSunkenThick:
		d.Line(x, y, x+w, y, wui.RGB(160, 160, 160))
		d.Line(x, y, x, y+h, wui.RGB(160, 160, 160))
		d.Line(x+w-1, y, x+w-1, y+h, wui.RGB(255, 255, 255))
		d.Line(x, y+h-1, x+w, y+h-1, wui.RGB(255, 255, 255))
		d.Line(x+1, y+1, x+w-1, y+1, wui.RGB(105, 105, 105))
		d.Line(x+1, y+1, x+1, y+h-1, wui.RGB(105, 105, 105))
		d.Line(x+w-2, y+1, x+w-2, y+h-1, wui.RGB(227, 227, 227))
		d.Line(x+1, y+h-2, x+w-1, y+h-2, wui.RGB(227, 227, 227))
	}
	innerX, innerY, _, _ := p.InnerBounds()
	drawContainer(p, makeOffsetDrawer(d, innerX, innerY))
}

func drawSlider(s *wui.Slider, d drawer) {
	var (
		drawSlideBar    func(offset int)
		drawCursorBody  func(offset, size int)
		drawCursorArrow func(offset int)
		// drawEndTicks and drawMiddleTicks are only assigned if ticks are
		// visible for this slider.
		drawEndTicks    = func(offset int) {}
		drawMiddleTicks = func(offset int) {}
	)

	cursorColor := wui.RGB(0, 120, 215)
	tickColor := wui.RGB(196, 196, 196)
	slideBarBorder := wui.RGB(214, 214, 214)
	slideBarBackground := wui.RGB(231, 231, 234)

	x, y, w, h := s.Bounds()
	if w <= 0 || h <= 0 {
		return
	}
	d.PushDrawRegion(x, y, w, h)
	defer d.PopDrawRegion()
	min, max := s.MinMax()
	innerTickCount := max - min - 1
	freq := s.TickFrequency()
	relCursor := s.CursorPosition() - min

	if s.Orientation() == wui.HorizontalSlider {
		xLeft := x + 13
		xRight := x + w - 14
		scale := 1.0 / float64(innerTickCount+1) * float64(xRight-xLeft)
		if xRight < xLeft {
			xRight = xLeft
			scale = 0
		}
		xOffset := int(float64(relCursor)*scale + 0.5)
		cursorCenter := xLeft + xOffset

		drawSlideBar = func(offset int) {
			if xLeft != xRight {
				d.DrawRect(x+8, y+offset, w-16, 4, slideBarBorder)
				d.FillRect(x+9, y+offset+1, w-18, 2, slideBarBackground)
			}
		}
		drawCursorBody = func(offset, size int) {
			d.FillRect(cursorCenter-5, y+offset, 11, size, cursorColor)
		}
		drawCursorArrow = func(offset int) {
			d.Polygon([]wui.Point{
				{X: int32(cursorCenter - 5), Y: int32(y + 15)},
				{X: int32(cursorCenter), Y: int32(y + 15 + offset)},
				{X: int32(cursorCenter + 5), Y: int32(y + 15)},
			}, cursorColor)
		}

		if s.TicksVisible() {
			drawEndTicks = func(offset int) {
				d.Line(xLeft, y+offset, xLeft, y+offset+4, tickColor)
				d.Line(xRight, y+offset, xRight, y+offset+4, tickColor)
			}
			drawMiddleTicks = func(offset int) {
				for i := freq; i <= innerTickCount; i += freq {
					x := xLeft + int(float64(i)*scale+0.5)
					d.Line(x, y+offset, x, y+offset+3, tickColor)
				}
			}
		}
	} else {
		yTop := y + 13
		yBottom := y + h - 14
		scale := 1.0 / float64(innerTickCount+1) * float64(yBottom-yTop)
		if yBottom < yTop {
			yBottom = yTop
			scale = 0
		}
		yOffset := int(float64(relCursor)*scale + 0.5)
		cursorCenter := yTop + yOffset

		drawSlideBar = func(offset int) {
			if yTop != yBottom {
				d.DrawRect(x+offset, y+8, 4, h-16, slideBarBorder)
				d.FillRect(x+offset+1, y+9, 2, h-18, slideBarBackground)
			}
		}
		drawCursorBody = func(offset, size int) {
			d.FillRect(x+offset, cursorCenter-5, size, 11, cursorColor)
		}
		drawCursorArrow = func(offset int) {
			d.Polygon([]wui.Point{
				{X: int32(x + 15), Y: int32(cursorCenter - 5)},
				{X: int32(x + 15 + offset), Y: int32(cursorCenter)},
				{X: int32(x + 15), Y: int32(cursorCenter + 5)},
			}, cursorColor)
		}

		if s.TicksVisible() {
			drawEndTicks = func(offset int) {
				d.Line(x+offset, yTop, x+offset+4, yTop, tickColor)
				d.Line(x+offset, yBottom, x+offset+4, yBottom, tickColor)
			}
			drawMiddleTicks = func(offset int) {
				for i := freq; i <= innerTickCount; i += freq {
					y := yTop + int(float64(i)*scale+0.5)
					d.Line(x+offset, y, x+offset+3, y, tickColor)
				}
			}
		}
	}

	switch s.TickPosition() {
	case wui.TicksBottomOrRight:
		drawSlideBar(8)
		drawCursorBody(2, 14)
		drawCursorArrow(5)
		drawEndTicks(22)
		drawMiddleTicks(22)
	case wui.TicksTopOrLeft:
		drawSlideBar(18)
		drawCursorBody(15, 14)
		drawCursorArrow(-5)
		drawEndTicks(5)
		drawMiddleTicks(6)
	case wui.TicksOnBothSides:
		drawSlideBar(19)
		drawCursorBody(10, 21)
		drawEndTicks(5)
		drawEndTicks(33)
		drawMiddleTicks(6)
		drawMiddleTicks(33)
	default:
		panic("unhandled tick position")
	}
}

func drawLabel(l *wui.Label, d drawer) {
	x, y, w, h := l.Bounds()
	textW, textH := d.TextExtent(l.Text())
	textX := x
	switch l.Alignment() {
	case wui.AlignCenter:
		textX = x + (w-textW)/2
	case wui.AlignRight:
		textX = x + w - textW
	}
	d.PushDrawRegion(x, y, w, h)
	d.TextOut(textX, y+(h-textH)/2, l.Text(), wui.RGB(0, 0, 0))
	d.PopDrawRegion()
}

func drawPaintBox(p *wui.PaintBox, d drawer) {
	x, y, w, h := p.Bounds()
	if w > 0 && h > 0 {
		d.DrawRect(x, y, w, h, wui.RGB(0, 0, 0))
		d.TextRectFormat(x, y, w, h, "Paint Box", wui.FormatCenter, wui.RGB(0, 0, 0))
	}
}

func drawIntUpDown(e *wui.IntUpDown, d drawer) {
	x, y, w, h := e.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		d.DrawRect(x, y, w, h, wui.RGB(122, 122, 122))
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(255, 255, 255))

		text := strconv.Itoa(e.Value())
		color := wui.RGB(0, 0, 0)
		d.TextOut(x+6, y+3, text, color)

		d.FillRect(x+w-19, y, 19, h, wui.RGB(231, 231, 231))
		d.DrawRect(x+w-19, y, 19, h, wui.RGB(172, 172, 172))
		d.DrawRect(x+w-19+2, y+2, 19-4, h-4, wui.RGB(172, 172, 172))
		d.DrawRect(x+w-19+2, y+h/2-1, 19-4, 2, wui.RGB(172, 172, 172))
		y1 := y + h/4
		d.Line(x+w-12, y1+2, x+w-12+5, y1+2, wui.RGB(0, 0, 0))
		d.Line(x+w-11, y1+1, x+w-11+3, y1+1, wui.RGB(0, 0, 0))
		d.Line(x+w-10, y1+0, x+w-10+1, y1+0, wui.RGB(0, 0, 0))
		y2 := y + 3*h/4 - 2
		d.Line(x+w-12, y2+0, x+w-12+5, y2+0, wui.RGB(0, 0, 0))
		d.Line(x+w-11, y2+1, x+w-11+3, y2+1, wui.RGB(0, 0, 0))
		d.Line(x+w-10, y2+2, x+w-10+1, y2+2, wui.RGB(0, 0, 0))
		d.PopDrawRegion()
	}
}

func drawFloatUpDown(e *wui.FloatUpDown, d drawer) {
	x, y, w, h := e.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		d.DrawRect(x, y, w, h, wui.RGB(122, 122, 122))
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(255, 255, 255))

		text := fmt.Sprintf("%."+strconv.Itoa(e.Precision())+"f", e.Value())
		color := wui.RGB(0, 0, 0)
		d.TextOut(x+6, y+3, text, color)

		d.FillRect(x+w-19, y, 19, h, wui.RGB(231, 231, 231))
		d.DrawRect(x+w-19, y, 19, h, wui.RGB(172, 172, 172))
		d.DrawRect(x+w-19+2, y+2, 19-4, h-4, wui.RGB(172, 172, 172))
		d.DrawRect(x+w-19+2, y+h/2-1, 19-4, 2, wui.RGB(172, 172, 172))
		y1 := y + h/4
		d.Line(x+w-12, y1+2, x+w-12+5, y1+2, wui.RGB(0, 0, 0))
		d.Line(x+w-11, y1+1, x+w-11+3, y1+1, wui.RGB(0, 0, 0))
		d.Line(x+w-10, y1+0, x+w-10+1, y1+0, wui.RGB(0, 0, 0))
		y2 := y + 3*h/4 - 2
		d.Line(x+w-12, y2+0, x+w-12+5, y2+0, wui.RGB(0, 0, 0))
		d.Line(x+w-11, y2+1, x+w-11+3, y2+1, wui.RGB(0, 0, 0))
		d.Line(x+w-10, y2+2, x+w-10+1, y2+2, wui.RGB(0, 0, 0))
		d.PopDrawRegion()
	}
}

func drawComboBox(c *wui.ComboBox, d drawer) {
	x, y, w, h := c.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		d.DrawRect(x, y, w, h, wui.RGB(173, 173, 173))
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(225, 225, 225))
		arrowX := x + w - 13
		arrowY := y + 9
		d.Line(arrowX, arrowY, arrowX+4, arrowY+4, wui.RGB(86, 86, 86))
		d.Line(arrowX+4, arrowY+3, arrowX+8, arrowY-1, wui.RGB(86, 86, 86))
		if w > 20 {
			i := c.SelectedIndex()
			items := c.Items()
			if 0 <= i && i < len(items) {
				text := items[i]
				d.PushDrawRegion(x, y, w-20, h)
				d.TextOut(x+4, y+4, text, wui.RGB(0, 0, 0))
				d.PopDrawRegion()
			}
		}
		d.PopDrawRegion()
	}
}

func drawProgressBar(p *wui.ProgressBar, d drawer) {
	x, y, w, h := p.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		d.DrawRect(x, y, w, h, wui.RGB(188, 188, 188))
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(230, 230, 230))
		if p.MovesForever() {
			if p.Vertical() {
				filledH := (h - 2) / 2
				d.FillRect(x+1, y+1+filledH/2, w-2, filledH, wui.RGB(0, 180, 40))
			} else {
				filledW := (w - 2) / 2
				d.FillRect(x+1+filledW/2, y+1, filledW, h-2, wui.RGB(0, 180, 40))
			}
		} else {
			if p.Vertical() {
				filledH := int(float64(h-2)*p.Value() + 0.5)
				d.FillRect(x+1, y+h-1-filledH, w-2, filledH, wui.RGB(0, 180, 40))
			} else {
				filledW := int(float64(w-2)*p.Value() + 0.5)
				d.FillRect(x+1, y+1, filledW, h-2, wui.RGB(0, 180, 40))
			}
		}
		d.PopDrawRegion()
	}
}

func drawEditLine(e *wui.EditLine, d drawer) {
	x, y, w, h := e.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		if e.Enabled() {
			d.DrawRect(x, y, w, h, wui.RGB(122, 122, 122))
		} else {
			d.DrawRect(x, y, w, h, wui.RGB(204, 204, 204))
		}
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(255, 255, 255))
		if e.ReadOnly() || !e.Enabled() {
			d.FillRect(x+2, y+2, w-4, h-4, wui.RGB(240, 240, 240))
		}
		text := e.Text()
		if e.IsPassword() {
			text = strings.Repeat("●", utf8.RuneCountInString(text))
		}
		color := wui.RGB(0, 0, 0)
		if !e.Enabled() {
			color = wui.RGB(109, 109, 109)
		}
		d.TextOut(x+6, y+3, text, color)
		d.PopDrawRegion()
	}
}

func drawTextEdit(t *wui.TextEdit, d drawer) {
	x, y, w, h := t.Bounds()
	if w > 0 && h > 0 {
		d.PushDrawRegion(x, y, w, h)
		if t.Enabled() {
			d.DrawRect(x, y, w, h, wui.RGB(122, 122, 122))
		} else {
			d.DrawRect(x, y, w, h, wui.RGB(204, 204, 204))
		}
		d.FillRect(x+1, y+1, w-2, h-2, wui.RGB(255, 255, 255))
		if !t.Enabled() {
			d.FillRect(x+2, y+2, w-4, h-4, wui.RGB(240, 240, 240))
		}
		color := wui.RGB(0, 0, 0)
		if !t.Enabled() {
			color = wui.RGB(109, 109, 109)
		}
		if t.WordWrap() {
			d.TextRectFormat(x+6, y+3, w-6, h-3, t.Text(), wui.FormatTopLeft, color)
		} else {
			d.TextOut(x+6, y+3, t.Text(), color)
		}
		d.PopDrawRegion()
	}
}
