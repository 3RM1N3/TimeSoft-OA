package main

import (
	"github.com/ying32/govcl/vcl"
)

//::private::
type TFormLoginFields struct {
}

func (f *TFormLogin) OnButtonLoginClick(sender vcl.IObject) {
	f.Close()
	vcl.Application.MainForm().Show()
}

func (f *TFormLogin) OnFormCreate(sender vcl.IObject) {

}

func (f *TFormLogin) OnLabel1Click(sender vcl.IObject) {

}
