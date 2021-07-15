package main

import (
	"os"
	"path"

	"fyne.io/fyne/v2/app"
)

func init() {
	wd, _ := os.Getwd()
	os.Setenv("FYNE_FONT", path.Join(wd, "simhei.ttf"))
}

func main() {
	nextPage := make(chan bool)
	thisApp := app.New()

	mainWindow := thisApp.NewWindow("时源办公自动化管理系统beta20210714-1")
	go func() {
		<-nextPage
		makeMainWindow(&mainWindow)
		mainWindow.Show()
	}()

	loginWindow := makeLoginPage(&thisApp, nextPage)
	loginWindow.ShowAndRun()
}
