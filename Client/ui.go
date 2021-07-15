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
	loginWindow := (*thisApp).NewWindow("登录")
	loginWindow.SetFixedSize(true)
	loginWindow.Resize(fyne.NewSize(300, 400))

	title := canvas.NewText("时源办公自动化管理系统", color.Black)
	title.Alignment = 1
	title.TextSize += 5

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("张三")

	serverIPAddr := widget.NewEntry()
	serverIPAddr.SetPlaceHolder("0.0.0.0:8080")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	signUp := widget.NewButton("注册账号", func() {
		fmt.Println("注册")
	})

	forgetPwd := widget.NewButton("忘记密码", func() {
		fmt.Println("忘记密码")
	})

	loginButton := widget.NewButton("登录", func() {})
	loginButton.OnTapped = func() {
		loginButton.Disable()

		name := nameEntry.Text
		result, err := login(serverIPAddr.Text, name, password.Text)
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
		nextPage <- true

		globalName = name
		globalID = result.ID

		loginWindow.Close()
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

func makeMainWindow(w *fyne.Window) {
	println("make main window")
	(*w).Resize(fyne.NewSize(600, 400))
	(*w).SetContent(widget.NewLabel("hello world"))
}
