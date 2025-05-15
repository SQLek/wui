package main

import "github.com/gonutz/wui/v2"

func main() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetTitle("Window")

	label1Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -30,
		Bold:   true,
	})

	label1 := wui.NewLabel()
	label1.SetFont(label1Font)
	label1.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndMax)
	label1.SetSize(100, 50)
	label1.SetText("UwU")
	label1.SetAlignment(wui.AlignCenter)
	window.Add(label1)

	window.Show()
}
