package main

import (
	"fmt"

	"github.com/ying32/govcl/vcl"
)

//::private::
type TFormLoginFields struct {
}

func (f *TLoginForm) OnLoginClick(sender vcl.IObject) {
	f.Close()
	vcl.Application.MainForm().Show()
}

func (f *TLoginForm) OnSignUpClick(sender vcl.IObject) {
	fmt.Println("注册账号")
}

func (f *TLoginForm) OnForgetPwdClick(sender vcl.IObject) {
	fmt.Println("忘记密码")
	//dlg := vcl.NewTaskDialog(f)
	//defer dlg.Free()
	//
	//dlg.SetCaption("提示")
	//dlg.SetText
	vcl.ShowMessage("请联系管理员重置密码")
}
