package main

import "github.com/ying32/govcl/vcl"

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&MainForm)
	vcl.Application.CreateForm(&LoginForm)
	//MainForm.Show()

	vcl.Application.SetShowMainForm(false)

	LoginForm.Show()

	vcl.Application.Run()
}
