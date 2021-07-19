package main

import "github.com/ying32/govcl/vcl"

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&MainForm)
	vcl.Application.CreateForm(&FormLogin)
	//MainForm.Show()

	vcl.Application.SetShowMainForm(false)

	FormLogin.Show()

	vcl.Application.Run()
}
