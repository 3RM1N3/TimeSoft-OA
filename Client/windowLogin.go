package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeLoginPage(thisApp *(fyne.App), nextPage chan bool) fyne.Window {
	loginWindow := (*thisApp).NewWindow("登录")  // 创建登录窗口
	loginWindow.SetFixedSize(true)             // 禁止改变窗口大小
	loginWindow.Resize(fyne.NewSize(300, 400)) // 设置大小300 400

	title := canvas.NewText("时源办公自动化管理系统", color.Black) // 设置logo标题
	title.Alignment = 1                                 // 居中对齐
	title.TextSize += 5                                 // 文字大小+5

	nameEntry := widget.NewEntry() // 账号输入框
	nameEntry.SetPlaceHolder("张三")

	serverIPAddr := widget.NewEntry() // 服务器ip地址输入框
	serverIPAddr.SetPlaceHolder("0.0.0.0:8080")

	password := widget.NewPasswordEntry() // 密码输入框
	password.SetPlaceHolder("Password")

	signUp := widget.NewButton("注册账号", func() {
		fmt.Println("注册")
	})

	forgetPwd := widget.NewButton("忘记密码", func() {
		fmt.Println("忘记密码")
	})

	loginButton := widget.NewButton("登录", func() {}) // 登录按钮
	loginButton.OnTapped = func() {
		loginButton.Disable() // 点击后禁用此按钮避免重复点击

		name := nameEntry.Text
		globalIPAddr = CompleteIPAddr(serverIPAddr.Text) // 设置全局ip地址

		result, err := postWithJsonLogin(globalIPAddr, name, password.Text)
		if err != nil {
			dialog.ShowInformation("错误", "无法连接服务器", loginWindow)
			loginButton.Enable()
			return
		}
		if result.Code != "200" {
			dialog.ShowInformation("错误", "账号或密码不正确", loginWindow)
			loginButton.Enable()
			return
		}
		nextPage <- true // 登录成功，可以切换页面

		globalName = name
		globalID = result.ID

		loginWindow.Close() // 关闭此窗口
	}

	rect := canvas.NewRectangle(color.White)

	loginWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		container.NewVBox(
			widget.NewLabel("服务器地址："),
			serverIPAddr,
			container.NewHBox(
				widget.NewLabel("账号："),
				layout.NewSpacer(),
				signUp,
			),
			nameEntry,
			container.NewHBox(
				widget.NewLabel("密码："),
				layout.NewSpacer(),
				forgetPwd,
			),
			password,
			rect,
			loginButton,
		),
	))

	return loginWindow
}
