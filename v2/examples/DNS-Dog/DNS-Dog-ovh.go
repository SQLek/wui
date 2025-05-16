package main

import "github.com/gonutz/wui/v2"

func OvhConfUi() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetInnerHeight(298)
	window.SetTitle("DNS Dog Configurator")
	window.SetHasMaxButton(false)
	window.SetResizable(false)

	panel1Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel1 := wui.NewPanel()
	panel1.SetFont(panel1Font)
	panel1.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndMax)
	panel1.SetBounds(5, 5, 589, 258)
	panel1.SetBorderStyle(wui.PanelBorderSingleLine)
	window.Add(panel1)

	ovhLbl := wui.NewLabel()
	ovhLbl.SetHorizontalAnchor(wui.AnchorMinAndMax)
	ovhLbl.SetBounds(-1, -2, 589, 26)
	ovhLbl.SetText("OVH connection configuration")
	ovhLbl.SetAlignment(wui.AlignCenter)
	panel1.Add(ovhLbl)

	endLbl := wui.NewLabel()
	endLbl.SetHorizontalAnchor(wui.AnchorMinAndMax)
	endLbl.SetBounds(10, 25, 52, 20)
	endLbl.SetText("Endpoint")
	panel1.Add(endLbl)

	editLine1 := wui.NewEditLine()
	editLine1.SetHorizontalAnchor(wui.AnchorMinAndMax)
	editLine1.SetBounds(70, 25, 120, 20)
	editLine1.SetText("ovh-eu")
	panel1.Add(editLine1)

	panel2Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel2 := wui.NewPanel()
	panel2.SetFont(panel2Font)
	panel2.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	panel2.SetBounds(4, 64, 580, 82)
	panel2.SetBorderStyle(wui.PanelBorderSingleLine)
	panel1.Add(panel2)

	clientLbl := wui.NewLabel()
	clientLbl.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	clientLbl.SetBounds(-1, -1, 580, 25)
	clientLbl.SetText("Client ID Configuration")
	clientLbl.SetAlignment(wui.AlignCenter)
	panel2.Add(clientLbl)

	clientidLbl := wui.NewLabel()
	clientidLbl.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	clientidLbl.SetBounds(5, 25, 72, 20)
	clientidLbl.SetText("Client ID:")
	panel2.Add(clientidLbl)

	clientScrtLbl := wui.NewLabel()
	clientScrtLbl.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	clientScrtLbl.SetBounds(5, 55, 77, 20)
	clientScrtLbl.SetText("Client secret:")
	panel2.Add(clientScrtLbl)

	clientIdInput := wui.NewEditLine()
	clientIdInput.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	clientIdInput.SetBounds(100, 25, 450, 20)
	panel2.Add(clientIdInput)

	clientScrtInput := wui.NewEditLine()
	clientScrtInput.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMinAndCenter)
	clientScrtInput.SetBounds(100, 55, 450, 20)
	clientScrtInput.SetIsPassword(true)
	panel2.Add(clientScrtInput)

	clientBtn := wui.NewRadioButton()
	clientBtn.SetHorizontalAnchor(wui.AnchorMax)
	clientBtn.SetBounds(264, 25, 100, 20)
	clientBtn.SetText("Use client ID")
	clientBtn.SetChecked(true)
	panel1.Add(clientBtn)

	appkeyBtn := wui.NewRadioButton()
	appkeyBtn.SetHorizontalAnchor(wui.AnchorMax)
	appkeyBtn.SetBounds(432, 25, 127, 20)
	appkeyBtn.SetText("Use application key")
	panel1.Add(appkeyBtn)

	panel3Font, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	panel3 := wui.NewPanel()
	panel3.SetFont(panel3Font)
	panel3.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	panel3.SetBounds(4, 150, 580, 102)
	panel3.SetBorderStyle(wui.PanelBorderSingleLine)
	panel1.Add(panel3)

	appkeyLbl := wui.NewLabel()
	appkeyLbl.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	appkeyLbl.SetX(-1)
	appkeyLbl.SetSize(580, 25)
	appkeyLbl.SetText("Application key Configuration")
	appkeyLbl.SetAlignment(wui.AlignCenter)
	panel3.Add(appkeyLbl)

	appkeyLbl1 := wui.NewLabel()
	appkeyLbl1.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	appkeyLbl1.SetBounds(5, 25, 89, 20)
	appkeyLbl1.SetText("Application Key:")
	panel3.Add(appkeyLbl1)

	appkeyLbl2 := wui.NewLabel()
	appkeyLbl2.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	appkeyLbl2.SetBounds(5, 50, 93, 20)
	appkeyLbl2.SetText("Application secret:")
	panel3.Add(appkeyLbl2)

	editLine4 := wui.NewEditLine()
	editLine4.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	editLine4.SetBounds(100, 25, 450, 20)
	editLine4.SetIsPassword(true)
	panel3.Add(editLine4)

	appkeyLbl3 := wui.NewLabel()
	appkeyLbl3.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	appkeyLbl3.SetBounds(5, 75, 86, 20)
	appkeyLbl3.SetText("Consumer key:")
	panel3.Add(appkeyLbl3)

	editLine5 := wui.NewEditLine()
	editLine5.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	editLine5.SetBounds(100, 50, 450, 20)
	editLine5.SetIsPassword(true)
	panel3.Add(editLine5)

	editLine6 := wui.NewEditLine()
	editLine6.SetAnchors(wui.AnchorMinAndMax, wui.AnchorMaxAndCenter)
	editLine6.SetBounds(100, 75, 450, 20)
	editLine6.SetIsPassword(true)
	panel3.Add(editLine6)

	saveBtn := wui.NewButton()
	saveBtn.SetVerticalAnchor(wui.AnchorMax)
	saveBtn.SetBounds(15, 270, 85, 25)
	saveBtn.SetText("Save")
	window.Add(saveBtn)

	cancelBtn := wui.NewButton()
	cancelBtn.SetAnchors(wui.AnchorMax, wui.AnchorMax)
	cancelBtn.SetBounds(500, 270, 85, 25)
	cancelBtn.SetText("Cancel")
	window.Add(cancelBtn)

	window.Show()
}
