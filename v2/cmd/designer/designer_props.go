package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gonutz/wui/v2"
)

const propMargin = 2

type uiProp struct {
	panel     *wui.Panel
	setter    string
	update    func()
	rightType func(t reflect.Type) bool
}

func (d *Designer) boolPanel(parent wui.Container, name string) (*wui.CheckBox, *wui.Panel) {
	c := wui.NewCheckBox()
	c.SetText(name)
	c.SetBounds(100, propMargin, 95, 17)
	p := wui.NewPanel()
	p.SetSize(195, c.Height()+2*propMargin)
	parent.Add(p)
	p.Add(c)
	return c, p
}

func (d *Designer) boolProp(name, getterFunc string) uiProp {
	c, p := d.boolPanel(d.w, name)
	setterFunc := "Set" + getterFunc // By convention.
	c.SetOnChange(func(on bool) {
		reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
			[]reflect.Value{reflect.ValueOf(on)},
		)
		d.updateProperties()
		d.preview.Paint()
	})
	update := func() {
		on := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0].Bool()
		if c.Checked() != on {
			c.SetChecked(on)
		}
	}
	rightType := func(t reflect.Type) bool {
		return t.Kind() == reflect.Bool
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}

func (d *Designer) intPanel(parent wui.Container, name string, minmax ...int) (*wui.IntUpDown, *wui.Panel) {
	n := wui.NewIntUpDown()
	n.SetOnTabFocus(n.SelectAll)
	if len(minmax) == 2 {
		n.SetMinMax(minmax[0], minmax[1])
	}
	n.SetBounds(100, propMargin, 90, 22)
	l := wui.NewLabel()
	l.SetText(name)
	l.SetAlignment(wui.AlignRight)
	// TODO This -1 might have to do with the below TODO about the IntUpDown
	// height.
	l.SetBounds(0, propMargin-1, 95, n.Height())
	p := wui.NewPanel()
	// TODO We add +2 to the height because for some reason setting the
	// height of an IntUpDown does not include the borders. Fix this in the
	// wui library.
	p.SetSize(195, n.Height()+2+2*propMargin)
	parent.Add(p)
	p.Add(l)
	p.Add(n)
	return n, p
}

func (d *Designer) intProp(name, getterFunc string, minmax ...int) uiProp {
	n, p := d.intPanel(d.w, name, minmax...)
	setterFunc := "Set" + getterFunc // By convention.
	n.SetOnValueChange(func(v int) {
		if d.active == nil {
			return
		}
		if m, ok := reflect.TypeOf(d.active).MethodByName(setterFunc); ok {
			reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
				[]reflect.Value{reflect.ValueOf(v).Convert(m.Type.In(1))},
			)
			d.updateProperties()
			d.preview.Paint()
		}
	})
	update := func() {
		v := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0]
		i := v.Convert(reflect.TypeOf(0)).Int()
		newValue := int(i)
		if n.Value() != newValue {
			n.SetValue(newValue)
		}
	}
	rightType := func(t reflect.Type) bool {
		return t.Kind() == reflect.Int || t.Kind() == reflect.Uint8
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}

func (d *Designer) floatProp(name, getterFunc string, minmax ...float64) uiProp {
	setterFunc := "Set" + getterFunc // By convention.
	n := wui.NewFloatUpDown()
	n.SetOnTabFocus(n.SelectAll)
	if len(minmax) == 2 {
		n.SetMinMax(minmax[0], minmax[1])
	}
	n.SetPrecision(6)
	n.SetBounds(100, propMargin, 90, 22)
	l := wui.NewLabel()
	l.SetText(name)
	l.SetAlignment(wui.AlignRight)
	// TODO This -1 might have to do with the below TODO about the
	// FloatUpDown height.
	l.SetBounds(0, propMargin-1, 95, n.Height())
	p := wui.NewPanel()
	// TODO We add +2 to the height because for some reason setting the
	// height of an FloatUpDown does not include the borders. Fix this in
	// the wui library.
	p.SetSize(195, n.Height()+2+2*propMargin)
	d.w.Add(p)
	p.Add(l)
	p.Add(n)
	n.SetOnValueChange(func(v float64) {
		if d.active == nil {
			return
		}
		if m, ok := reflect.TypeOf(d.active).MethodByName(setterFunc); ok {
			reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
				[]reflect.Value{reflect.ValueOf(v).Convert(m.Type.In(1))},
			)
			d.updateProperties()
			d.preview.Paint()
		}
	})
	update := func() {
		v := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0]
		newValue := v.Convert(reflect.TypeOf(0.0)).Float()
		if n.Value() != newValue {
			n.SetValue(newValue)
		}
	}
	rightType := func(t reflect.Type) bool {
		return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}

func (d *Designer) stringPanel(parent wui.Container, name string) (*wui.EditLine, *wui.Panel) {
	t := wui.NewEditLine()
	t.SetOnTabFocus(t.SelectAll)
	t.SetBounds(100, propMargin, 90, 22)
	l := wui.NewLabel()
	l.SetText(name)
	l.SetAlignment(wui.AlignRight)
	l.SetBounds(0, propMargin-1, 95, t.Height())
	p := wui.NewPanel()
	p.SetSize(195, t.Height()+2*propMargin)
	parent.Add(p)
	p.Add(l)
	p.Add(t)
	return t, p
}

func (d *Designer) stringProp(name, getterFunc string) uiProp {
	t, p := d.stringPanel(d.w, name)
	setterFunc := "Set" + getterFunc // By convention.
	t.SetOnTextChange(func() {
		if d.active == nil {
			return
		}
		if _, ok := reflect.TypeOf(d.active).MethodByName(setterFunc); ok {
			reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
				[]reflect.Value{reflect.ValueOf(t.Text())},
			)
			d.updateProperties()
			d.preview.Paint()
		}
	})
	update := func() {
		text := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0].String()
		if t.Text() != text {
			t.SetText(text)
		}
	}
	rightType := func(t reflect.Type) bool {
		return t.Kind() == reflect.String
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}

func (d *Designer) stringListProp(name, getterFunc string) uiProp {
	setterFunc := "Set" + getterFunc // By convention.
	l := wui.NewLabel()
	l.SetBounds(10, 5, 180, 13)
	l.SetText(name)
	l.SetAlignment(wui.AlignCenter)
	list := wui.NewTextEdit()
	list.SetBounds(10, 20, 180, 80)
	list.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndMax)
	p := wui.NewPanel()
	p.SetSize(195, list.Height()+2*propMargin)
	d.w.Add(p)
	p.Add(l)
	p.Add(list)
	list.SetOnTextChange(func() {
		if d.active == nil {
			return
		}
		if _, ok := reflect.TypeOf(d.active).MethodByName(setterFunc); ok {
			items := strings.Split(list.Text(), "\r\n")
			items = removeEmptyStrings(items)
			l.SetText(fmt.Sprintf("%s (%d)", name, len(items)))
			reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
				[]reflect.Value{reflect.ValueOf(items)},
			)
			start, end := list.CursorPosition()
			d.updateProperties()
			list.SetSelection(start, end)
			d.preview.Paint()
		}
	})
	update := func() {
		items := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0].Interface().([]string)
		l.SetText(fmt.Sprintf("%s (%d)", name, len(items)))
		newText := strings.Join(items, "\r\n") + "\r\n"
		if list.Text() != newText {
			list.SetText(newText)
		}
	}
	rightType := func(t reflect.Type) bool {
		// NOTE that currently there is only []string, we might have to
		// check for the underlying slice type if we support others in the
		// future.
		return t.Kind() == reflect.Slice
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}

// enumNames must correspond to the respective const, the order is important
// and the consts must be iota'd, i.e. start with 0 and increment by 1.
func (d *Designer) enumProp(name, getterFunc string, enumNames ...string) uiProp {
	setterFunc := "Set" + getterFunc // By convention.
	c := wui.NewComboBox()
	for _, name := range enumNames {
		c.AddItem(name)
	}
	c.SetBounds(100, propMargin, 90, 22)
	l := wui.NewLabel()
	l.SetText(name)
	l.SetAlignment(wui.AlignRight)
	l.SetBounds(0, propMargin-1, 95, c.Height())
	p := wui.NewPanel()
	p.SetSize(195, c.Height()+2*propMargin)
	d.w.Add(p)
	p.Add(l)
	p.Add(c)
	c.SetOnChange(func(index int) {
		m, ok := reflect.TypeOf(d.active).MethodByName(setterFunc)
		if ok {
			reflect.ValueOf(d.active).MethodByName(setterFunc).Call(
				[]reflect.Value{reflect.ValueOf(index).Convert(m.Type.In(1))},
			)
			d.updateProperties()
			d.preview.Paint()
		}
	})
	update := func() {
		v := reflect.ValueOf(d.active).MethodByName(getterFunc).Call(nil)[0]
		index := int(v.Convert(reflect.TypeOf(0)).Int())
		if c.SelectedIndex() != index {
			c.SetSelectedIndex(index)
		}
	}
	rightType := func(t reflect.Type) bool {
		return true
	}
	return uiProp{
		panel:     p,
		setter:    setterFunc,
		update:    update,
		rightType: rightType,
	}
}
