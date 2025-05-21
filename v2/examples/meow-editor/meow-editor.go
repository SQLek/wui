package main

import "github.com/gonutz/wui/v2"

func main() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerY(51)
	window.SetInnerSize(644, 576)
	window.SetTitle("Meow Editor")

	TitlePnlFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	TitlePnl := wui.NewPanel()
	TitlePnl.SetFont(TitlePnlFont)
	TitlePnl.SetHorizontalAnchor(wui.AnchorMinAndMax)
	TitlePnl.SetBounds(18, 10, 608, 133)
	TitlePnl.SetBorderStyle(wui.PanelBorderSunken)
	window.Add(TitlePnl)

	TitleLabelFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -35,
		Bold:   true,
	})

	TitleLabel := wui.NewLabel()
	TitleLabel.SetFont(TitleLabelFont)
	TitleLabel.SetHorizontalAnchor(wui.AnchorMinAndMax)
	TitleLabel.SetSize(608, 133)
	TitleLabel.SetText("Say Meow")
	TitleLabel.SetAlignment(wui.AlignCenter)
	TitlePnl.Add(TitleLabel)

	textInput := wui.NewTextEdit()
	textInput.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndMax)
	textInput.SetBounds(19, 154, 606, 347)
	textInput.SetText("Type it here")
	textInput.SetWordWrap(true)
	window.Add(textInput)

	NewBtn := wui.NewButton()
	NewBtn.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
	NewBtn.SetBounds(20, 515, 130, 50)
	NewBtn.SetText("New")
	window.Add(NewBtn)

	LoadBtn := wui.NewButton()
	LoadBtn.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
	LoadBtn.SetBounds(250, 515, 130, 50)
	LoadBtn.SetText("Load")
	window.Add(LoadBtn)

	SaveBtn := wui.NewButton()
	SaveBtn.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
	SaveBtn.SetBounds(490, 515, 130, 50)
	SaveBtn.SetText("Save")
	window.Add(SaveBtn)

	window.Show()
}
