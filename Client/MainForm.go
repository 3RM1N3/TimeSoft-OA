package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"sort"
)

// 提交后显示今日工作量

type TMainForm struct { // 主窗体
	*vcl.TForm // 继承窗体类
	PageControl *vcl.TPageControl // 页面控制器
	StatusBar   *vcl.TStatusBar // 底部状态栏

	// 扫描页面
	PageScan      *vcl.TTabSheet
	BtnSelectDir  *vcl.TButton
	BtnSubmitScan *vcl.TButton
	CoComboBox    *vcl.TComboBox
	ListView      *vcl.TListView

	// 修图页面
	PageEditPic   *vcl.TTabSheet
	BtnGetMission *vcl.TButton
	BtnSubmitEdit *vcl.TButton
	EditTaskList  *vcl.TListView

	// 返工页面
	PageRework   *vcl.TTabSheet
	BtnGetRework *vcl.TButton
	BtnSubmitRwk *vcl.TButton
	RwkTaskList  *vcl.TListView
}

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	// 主窗口
	f.SetCaption("时源科技")
	constraints := f.Constraints()
	constraints.SetMinHeight(600)
	constraints.SetMinWidth(860)
	f.SetConstraints(constraints)
	f.SetHeight(600)
	f.SetPosition(types.PoDesktopCenter)
	f.SetWidth(860)
	f.SetDoubleBuffered(true) // 开启双缓冲

	// 状态栏
	f.StatusBar = vcl.NewStatusBar(f)
	f.StatusBar.SetParent(f)
	f.StatusBar.SetSimpleText("就绪")

	// 标签页面控制器
	f.PageControl = vcl.NewPageControl(f)
	f.PageControl.SetParent(f)
	f.PageControl.SetActivePageIndex(0)
	f.PageControl.SetAnchors(types.NewSet( // 设置锚点
		types.AkBottom, types.AkLeft, types.AkRight, types.AkTop))
	f.PageControl.SetHeight(574)
	f.PageControl.SetWidth(860)
	f.PageControl.SetOnChanging(func(sender vcl.IObject, allowChange *bool) {
		if MissionInProgress {
			vcl.ShowMessageFmt("当前尚有任务未提交！")
		}
		*allowChange = !MissionInProgress
	})

	// 三个标签
	f.MakeScanTab()    // 扫描页面
	f.MakeEditPicTab() // 修图页面
	f.MakeReworkTab()  // 返工页面

}

// 窗口显示时进行的操作
func (f *TMainForm) OnFormShow() {
	// 获取本机第一个mac地址
	macs, err := GetMacAddrs()
	if err != nil {
		vcl.ShowMessageFmt("获取本机MAC地址失败：%v", err)
		f.Close()
	}
	sort.Strings(macs)
	macAddr = macs[0]

	missionNum, err := getEditMission() // 获取修图页面的表格
	if err != nil {
		vcl.ShowMessageFmt("获取修图任务错误：%v", err)
		return
	}

	if missionNum == 0 {
		MainForm.BtnGetMission.SetEnabled(false) // 禁用按钮
	} else {
		vcl.ShowMessageFmt("查询到%d个修图任务，请及时处理。", missionNum)
	}

	// 获取并设置现有客户公司
	clientcos, err := GetClientCo()
	if err != nil {
		vcl.ShowMessageFmt("获取已有客户公司时出现错误：%v", err)
		f.Close()
	}
	items := MainForm.CoComboBox.Items()
	items.AddStrings2(clientcos) // 设置下拉列表的值
}
