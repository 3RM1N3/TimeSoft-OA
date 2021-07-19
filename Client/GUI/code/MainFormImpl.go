package main

import (
	"github.com/ying32/govcl/vcl"
)

//::private::
type TMainFormFields struct {
}

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.ListView1.SetRowSelect(true)
	f.ListView1.SetReadOnly(true)
	f.ListView1.SetGridLines(true)
}

func (f *TMainForm) OnPageControl1Change(sender vcl.IObject) {

}

func (f *TMainForm) RefreshPageScan() {

}
