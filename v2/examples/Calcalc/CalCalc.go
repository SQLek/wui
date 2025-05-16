package main

import "github.com/gonutz/wui/v2"

func main() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerSize(337, 145)
	window.SetTitle("CalCalc")

	ExtrtValFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	ExtrtVal := wui.NewPanel()
	ExtrtVal.SetFont(ExtrtValFont)
	ExtrtVal.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorMinAndCenter)
	ExtrtVal.SetBounds(5, 4, 160, 50)
	ExtrtVal.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(ExtrtVal)

	ExtrLngthIn := wui.NewEditLine()
	ExtrLngthIn.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorMinAndCenter)
	ExtrLngthIn.SetBounds(5, 20, 150, 20)
	ExtrtVal.Add(ExtrLngthIn)

	ExtAmmLbl := wui.NewLabel()
	ExtAmmLbl.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorMinAndCenter)
	ExtAmmLbl.SetBounds(5, 3, 150, 13)
	ExtAmmLbl.SetText("Extrude length")
	ExtrtVal.Add(ExtAmmLbl)

	ActRtPnlFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	ActRtPnl := wui.NewPanel()
	ActRtPnl.SetFont(ActRtPnlFont)
	ActRtPnl.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorMinAndCenter)
	ActRtPnl.SetBounds(170, 4, 160, 50)
	ActRtPnl.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(ActRtPnl)

	ActRtDstLbl := wui.NewLabel()
	ActRtDstLbl.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorCenter)
	ActRtDstLbl.SetBounds(5, 3, 150, 13)
	ActRtDstLbl.SetText("Actual Rotation distance")
	ActRtPnl.Add(ActRtDstLbl)

	ActRtDstIn := wui.NewEditLine()
	ActRtDstIn.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorCenter)
	ActRtDstIn.SetBounds(5, 20, 150, 20)
	ActRtPnl.Add(ActRtDstIn)

	panel2Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel2 := wui.NewPanel()
	panel2.SetFont(panel2Font)
	panel2.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorMaxAndCenter)
	panel2.SetBounds(5, 60, 160, 50)
	panel2.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(panel2)

	DsrdExtLbl := wui.NewLabel()
	DsrdExtLbl.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorCenter)
	DsrdExtLbl.SetBounds(5, 3, 150, 13)
	DsrdExtLbl.SetText("Actual extrude length")
	panel2.Add(DsrdExtLbl)

	ActExtIn := wui.NewEditLine()
	ActExtIn.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorCenter)
	ActExtIn.SetBounds(5, 20, 150, 20)
	panel2.Add(ActExtIn)

	panel1Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel1 := wui.NewPanel()
	panel1.SetFont(panel1Font)
	panel1.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorMaxAndCenter)
	panel1.SetBounds(170, 60, 160, 50)
	panel1.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(panel1)

	NewRtDst := wui.NewLabel()
	NewRtDst.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorCenter)
	NewRtDst.SetBounds(5, 3, 150, 13)
	NewRtDst.SetText("Calculated rotation distance")
	panel1.Add(NewRtDst)

	CalcRtDstFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
		Bold:   true,
	})

	CalcRtDst := wui.NewEditLine()
	CalcRtDst.SetFont(CalcRtDstFont)
	CalcRtDst.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorCenter)
	CalcRtDst.SetBounds(5, 20, 150, 20)
	CalcRtDst.SetText("0")
	CalcRtDst.SetReadOnly(true)
	panel1.Add(CalcRtDst)

	resetBtn := wui.NewButton()
	resetBtn.SetVerticalAnchor(wui.AnchorMax)
	resetBtn.SetBounds(5, 115, 80, 25)
	resetBtn.SetText("Reset")
	window.Add(resetBtn)

	InfoBtn := wui.NewButton()
	InfoBtn.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
	InfoBtn.SetBounds(85, 115, 80, 25)
	InfoBtn.SetText("Info")
	InfoBtn.SetOnClick(func() {
		infoWindow()
	})
	window.Add(InfoBtn)

	AboutBtn := wui.NewButton()
	AboutBtn.SetAnchors(wui.AnchorCenter, wui.AnchorMax)
	AboutBtn.SetBounds(170, 115, 80, 25)
	AboutBtn.SetText("About")
	window.Add(AboutBtn)

	ExitBtn := wui.NewButton()
	ExitBtn.SetAnchors(wui.AnchorMax, wui.AnchorMax)
	ExitBtn.SetBounds(250, 115, 80, 25)
	ExitBtn.SetText("Exit")
	window.Add(ExitBtn)

	window.Show()
}
