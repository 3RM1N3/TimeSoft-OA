package main

import (
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func MakeMainWindow(w *fyne.Window) {
	(*w).Resize(fyne.NewSize(600, 400))
	(*w).SetFixedSize(true)

	(*w).SetContent(
		container.NewAppTabs(
			container.NewTabItem("扫描", makeScanPage(w)),
			container.NewTabItem("修图", makecEditImagePage()),
			container.NewTabItem("待返工", makeToBeReworkedPage()),
		),
	)
}

// 扫描页面
func makeScanPage(w *fyne.Window) fyne.CanvasObject {
	// 获取之前完成的项目
	scanedList := [][]string{
		{"张三丰0715", "15", "174", "2021-07-15 15:54"},
		{"张三丰0713", "10", "200", "2021-07-13 15:54"},
		{"张三丰0712", "11", "123", "2021-07-12 15:54"},
	}

	uploadButton := widget.NewButton("上传", func() {}) // 提交上传按钮
	uploadButton.OnTapped = func() {
		uploadButton.Disable()
		defer uploadButton.Enable()

		s := FolderChooser()
		if s == "" {
			return
		}

		if err := Zip(s, path.Base(s)+".zip"); err != nil {
			dialog.ShowInformation("错误", "文件夹压缩失败，请稍候或修改文件夹名后重试", *w)
		}
	}

	return container.NewBorder(
		container.NewVBox(
			uploadButton,
			widget.NewSeparator(),
			widget.NewLabel("已完成："),
		),
		nil,
		nil,
		nil,
		makeTable(scanedList),
	)
}

func makeTable(scanedList [][]string) fyne.CanvasObject {
	scanedList = append([][]string{{"文件夹名", "子文件夹数", "总文件数", "上传时间"}}, scanedList...)

	listLength := len(scanedList)
	t := widget.NewTable(
		func() (int, int) {
			return listLength, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("-")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Alignment = 1

			if id.Col == 0 {
				if id.Row == 0 {
					label.SetText("序号")
					return
				}
				label.SetText(fmt.Sprintf("%d", listLength-id.Row))
				return
			}
			label.SetText(scanedList[id.Row][id.Col-1])
		})

	t.SetColumnWidth(0, 40)
	t.SetColumnWidth(1, 140)
	t.SetColumnWidth(2, 85)
	t.SetColumnWidth(3, 85)
	t.SetColumnWidth(4, 190)

	return t
}

// 修图页面
func makecEditImagePage() fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("修图页面"))
}

// 错误处理页面
func makeToBeReworkedPage() fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("待返工页面"))
}
