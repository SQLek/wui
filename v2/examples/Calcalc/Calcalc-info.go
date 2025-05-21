package main

import "github.com/gonutz/wui/v2"

func infoWindow() {
	InfoWinFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	InfoWin := wui.NewWindow()
	InfoWin.SetFont(InfoWinFont)
	InfoWin.SetInnerSize(599, 401)
	InfoWin.SetTitle("Info")
	InfoWin.SetResizable(false)

	panel1Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel1 := wui.NewPanel()
	panel1.SetFont(panel1Font)
	panel1.SetBounds(4, 4, 590, 345)
	panel1.SetBorderStyle(wui.PanelBorderSingleLine)
	InfoWin.Add(panel1)

	textEdit1 := wui.NewTextEdit()
	textEdit1.SetBounds(5, 4, 579, 335)
	textEdit1.SetText("your info text")
	textEdit1.SetWordWrap(true)
	panel1.Add(textEdit1)

	CloseBtn := wui.NewButton()
	CloseBtn.SetBounds(240, 363, 85, 25)
	CloseBtn.SetText("Close")
	InfoWin.Add(CloseBtn)

	InfoWin.Show()
}
