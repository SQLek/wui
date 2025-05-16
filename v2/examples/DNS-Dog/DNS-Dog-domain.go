package main

import "github.com/gonutz/wui/v2"

func DomainConfUi() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerHeight(236)
	window.SetTitle("DNS Dog Configurator")
	window.SetHasMaxButton(false)
	window.SetResizable(false)

	panel1Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel1 := wui.NewPanel()
	panel1.SetFont(panel1Font)
	panel1.SetBounds(5, 5, 590, 193)
	panel1.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(panel1)

	dmnLbl := wui.NewLabel()
	dmnLbl.SetSize(590, 20)
	dmnLbl.SetText("Domains configuration")
	dmnLbl.SetAlignment(wui.AlignCenter)
	panel1.Add(dmnLbl)

	panel2Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel2 := wui.NewPanel()
	panel2.SetFont(panel2Font)
	panel2.SetBounds(5, 20, 580, 81)
	panel2.SetBorderStyle(wui.PanelBorderSingleLine)
	panel1.Add(panel2)

	zoneLbl := wui.NewLabel()
	zoneLbl.SetSize(580, 20)
	zoneLbl.SetText("Zone")
	zoneLbl.SetAlignment(wui.AlignCenter)
	panel2.Add(zoneLbl)

	zoneBox := wui.NewComboBox()
	zoneBox.SetBounds(140, 20, 300, 25)
	zoneBox.SetItems([]string{"Combo Box"})
	zoneBox.SetSelectedIndex(0)
	panel2.Add(zoneBox)

	zoneLbl2 := wui.NewLabel()
	zoneLbl2.SetBounds(5, 20, 124, 25)
	zoneLbl2.SetText("Zones:")
	panel2.Add(zoneLbl2)

	zoneDel := wui.NewButton()
	zoneDel.SetBounds(465, 20, 85, 25)
	zoneDel.SetText("Delete")
	panel2.Add(zoneDel)

	zoneTxt := wui.NewEditLine()
	zoneTxt.SetBounds(140, 55, 300, 20)
	zoneTxt.SetText("domain.com")
	panel2.Add(zoneTxt)

	zoneTextLbl := wui.NewLabel()
	zoneTextLbl.SetBounds(5, 55, 150, 20)
	zoneTextLbl.SetText("Selected zone:")
	panel2.Add(zoneTextLbl)

	zoneAddBtn1 := wui.NewButton()
	zoneAddBtn1.SetBounds(465, 50, 85, 25)
	zoneAddBtn1.SetText("Add/Edit")
	panel2.Add(zoneAddBtn1)

	panel3Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel3 := wui.NewPanel()
	panel3.SetFont(panel3Font)
	panel3.SetBounds(5, 105, 580, 81)
	panel3.SetBorderStyle(wui.PanelBorderSingleLine)
	panel1.Add(panel3)

	subDomLbl := wui.NewLabel()
	subDomLbl.SetSize(580, 20)
	subDomLbl.SetText("Subdomains")
	subDomLbl.SetAlignment(wui.AlignCenter)
	panel3.Add(subDomLbl)

	subdomLbl := wui.NewLabel()
	subdomLbl.SetBounds(5, 20, 100, 20)
	subdomLbl.SetText("Subdomain:")
	panel3.Add(subdomLbl)

	subdomLbl2 := wui.NewLabel()
	subdomLbl2.SetBounds(5, 55, 100, 20)
	subdomLbl2.SetText("Selected subdomain:")
	panel3.Add(subdomLbl2)

	subdomBox := wui.NewComboBox()
	subdomBox.SetBounds(140, 20, 300, 21)
	subdomBox.SetItems([]string{"Combo Box"})
	subdomBox.SetSelectedIndex(0)
	panel3.Add(subdomBox)

	subdomDelBtn := wui.NewButton()
	subdomDelBtn.SetBounds(465, 20, 85, 25)
	subdomDelBtn.SetText("Delete")
	panel3.Add(subdomDelBtn)

	subdomAddBtn := wui.NewButton()
	subdomAddBtn.SetBounds(465, 50, 85, 25)
	subdomAddBtn.SetText("Add/Edit")
	panel3.Add(subdomAddBtn)

	subdomTxt := wui.NewEditLine()
	subdomTxt.SetBounds(140, 55, 300, 20)
	subdomTxt.SetText("subdomain")
	panel3.Add(subdomTxt)

	saveBtn := wui.NewButton()
	saveBtn.SetBounds(15, 205, 85, 25)
	saveBtn.SetText("Save")
	window.Add(saveBtn)

	cancelBtn := wui.NewButton()
	cancelBtn.SetBounds(500, 205, 85, 25)
	cancelBtn.SetText("Cancel")
	window.Add(cancelBtn)

	window.Show()
}
