package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

// 创建页面

// 创建扫描标签页
func (f *TMainForm) MakeScanTab() {
	f.PageScan = vcl.NewTabSheet(f)
	f.PageScan.SetParent(f.PageControl)
	f.PageScan.SetCaption("扫描")
	f.PageScan.SetHeight(544)
	f.PageScan.SetWidth(852)

	// 选择文件夹按钮
	f.BtnSelectDir = vcl.NewButton(f)
	f.BtnSelectDir.SetParent(f.PageScan)
	f.BtnSelectDir.SetCaption("选择工作文件夹")
	f.BtnSelectDir.SetOnClick(f.OnScanSelectDir)
	f.BtnSelectDir.SetHeight(30)
	f.BtnSelectDir.SetWidth(100)
	f.BtnSelectDir.SetLeft(10)
	f.BtnSelectDir.SetTop(10)

	// 提交按钮
	f.BtnSubmitScan = vcl.NewButton(f)
	f.BtnSubmitScan.SetParent(f.PageScan)
	f.BtnSubmitScan.SetCaption("提交")
	f.BtnSubmitScan.SetOnClick(f.OnScanSubmit)
	f.BtnSubmitScan.SetHeight(30)
	f.BtnSubmitScan.SetWidth(60)
	f.BtnSubmitScan.SetLeft(160)
	f.BtnSubmitScan.SetTop(10)

	// 选择已有项目公司的下拉列表
	f.CoComboBox = vcl.NewComboBox(f)
	f.CoComboBox.SetParent(f.PageScan)
	f.CoComboBox.SetTop(60)
	f.CoComboBox.SetLeft(10)
	f.CoComboBox.DoubleBuffered()

	// 列表
	f.ListView = vcl.NewListView(f)
	f.ListView.SetParent(f.PageScan)
	f.ListView.SetAnchors(types.NewSet( // 设置锚点
		types.AkBottom, types.AkLeft, types.AkRight, types.AkTop))
	f.ListView.SetGridLines(true)
	f.ListView.SetHeight(456)
	f.ListView.SetReadOnly(true)
	f.ListView.SetRowSelect(true)
	f.ListView.SetTop(88)
	f.ListView.SetViewStyle(types.VsReport)
	f.ListView.SetWidth(848)

	col := f.ListView.Columns().Add()
	col.SetCaption("序号")
	col.SetWidth(48)

	col = f.ListView.Columns().Add()
	col.SetCaption("档号")
	col.SetWidth(695)

	col = f.ListView.Columns().Add()
	col.SetCaption("文件数")
	col.SetWidth(100)
}

// 创建修图标签页
func (f *TMainForm) MakeEditPicTab() {
	f.PageEditPic = vcl.NewTabSheet(f)
	f.PageEditPic.SetParent(f.PageControl)
	f.PageEditPic.SetCaption("修图")
	f.PageEditPic.SetHeight(544)
	f.PageEditPic.SetWidth(852)

	// 获取任务按钮
	f.BtnGetMission = vcl.NewButton(f)
	f.BtnGetMission.SetParent(f.PageEditPic)
	f.BtnGetMission.SetCaption("获取任务")
	f.BtnGetMission.SetOnClick(f.OnEditGetMission)
	f.BtnGetMission.SetHeight(30)
	f.BtnGetMission.SetWidth(100)
	f.BtnGetMission.SetLeft(10)
	f.BtnGetMission.SetTop(10)

	// 提交按钮
	f.BtnSubmitEdit = vcl.NewButton(f)
	f.BtnSubmitEdit.SetParent(f.PageEditPic)
	f.BtnSubmitEdit.SetCaption("提交")
	f.BtnSubmitEdit.SetOnClick(f.OnEditSubmit)
	f.BtnSubmitEdit.SetHeight(30)
	f.BtnSubmitEdit.SetWidth(60)
	f.BtnSubmitEdit.SetLeft(160)
	f.BtnSubmitEdit.SetTop(10)

	// 列表
	f.EditTaskList = vcl.NewListView(f)
	f.EditTaskList.SetParent(f.PageEditPic)
	f.EditTaskList.SetAnchors(types.NewSet( // 设置锚点
		types.AkBottom, types.AkLeft, types.AkRight, types.AkTop))
	f.EditTaskList.SetGridLines(true)
	f.EditTaskList.SetHeight(456)
	f.EditTaskList.SetReadOnly(true)
	f.EditTaskList.SetRowSelect(true)
	f.EditTaskList.SetTop(88)
	f.EditTaskList.SetViewStyle(types.VsReport)
	f.EditTaskList.SetWidth(848)
	f.EditTaskList.SetOnDblClick(f.SetRework) // 双击设置返工

	col := f.EditTaskList.Columns().Add()
	col.SetCaption("序号")
	col.SetWidth(48)

	col = f.EditTaskList.Columns().Add()
	col.SetCaption("档号")
	col.SetWidth(695)
}

// 创建返工标签页
func (f *TMainForm) MakeReworkTab() {
	f.PageRework = vcl.NewTabSheet(f)
	f.PageRework.SetParent(f.PageControl)
	f.PageRework.SetCaption("返工")
	f.PageRework.SetHeight(544)
	f.PageRework.SetWidth(852)
}

// 相关方法

// 扫描页面

// 选择扫描项目文件夹按钮
func (f *TMainForm) OnScanSelectDir(sender vcl.IObject) {
	if ok, dir := vcl.SelectDirectory2("选择工作文件夹", "C:/", true); ok {
		err := SetProjectDir(dir)
		if err != nil {
			vcl.ShowMessageFmt("文件夹校验失败，请选择曾用的项目文件夹或空文件夹：%v", err)
			return
		}
		ProjectDir = dir
		ChStartScan <- true
		MissionInProgress = true
		f.BtnSelectDir.SetEnabled(false)
	}
}

// 扫描提交按钮
func (f *TMainForm) OnScanSubmit(sender vcl.IObject) {
	// 打包提交
	clientCo := f.CoComboBox.Text()
	if clientCo == "" {
		vcl.ShowMessageFmt("尚未选择客户公司")
		return
	}

	OverScan = true
	s := "尚未设置项目文件夹\n"
	if ProjectDir != "" {
		err := PackSubmitDir(ProjectDir, clientCo, globalPhone, 0)
		if err != nil {
			vcl.ShowMessageFmt("提交错误：%v", err)
			return
		}
		s = "提交成功！\n"
	}
	ProjectDir = ""
	f.StatusBar.SetSimpleText("就绪")
	workload, err := TodayWorkload()
	if err != nil {
		vcl.ShowMessageFmt("%s获取今日已提交任务数失败：%v", s, err)
		return
	}
	MissionInProgress = false
	f.BtnSelectDir.SetEnabled(true)
	vcl.ShowMessageFmt("%s今日任务量：\n扫描：%d\n修图：%d\n返工：%d", s, workload.Scan, workload.Edit, workload.Rework)
}

// 修图页面

// 设置修图页面的列表
func (f *TMainForm) SetEditList(l []string) {
	MainForm.EditTaskList.Items().BeginUpdate()

	var lv1 *vcl.TListItem
	for i, j := range l {
		lv1 = MainForm.EditTaskList.Items().Add()
		lv1.SetCaption(strconv.Itoa(i + 1))
		lv1.SubItems().Add(j)
	}

	MainForm.EditTaskList.Items().EndUpdate()
}

// 双击条目返工
func (f *TMainForm) SetRework(sender vcl.IObject) {
	i := vcl.MessageDlg("确定要将xx标记为返工吗？", types.MtConfirmation, types.MbYes, types.MbNo)
	if i == 7 { // 不返工此任务
		return
	}

	// 返工此任务
	selString := f.EditTaskList.Selected().SubItems().Strings(0)
	fmt.Println(selString)

	err := ReworkItem(selString)
	if err != nil {
		vcl.ShowMessageFmt("无法返工此任务：%v", err)
		return
	}

	err = os.RemoveAll(filepath.Join(ProjectDir, selString))
	if err != nil {
		vcl.ShowMessageFmt("删除本地文件失败：%v", err)
		return
	}

	f.EditTaskList.DeleteSelected() // 从列表中删除条目
}

// 修图页面获取任务按钮
func (f *TMainForm) OnEditGetMission(sender vcl.IObject) {

	ok, dir := vcl.SelectDirectory2("选择任务文件保存位置", "C:/", true)
	if !ok {
		log.Println("未选择文件夹")
		return
	}
	ProjectDir = dir
	log.Println("选择文件夹：", ProjectDir)

	fileIDs, err := SaveAndUnzip(ProjectDir)
	if err != nil {
		vcl.ShowMessageFmt("获取任务时出现错误：%v", err)
		return
	}

	f.SetEditList(fileIDs)

	MissionInProgress = true
	f.BtnGetMission.SetEnabled(false)
}

// 修图提交按钮
func (f *TMainForm) OnEditSubmit(sender vcl.IObject) {
	// 打包提交
	s := "尚未设置项目文件夹\n"
	if ProjectDir != "" {
		err := PackSubmitDir(ProjectDir, "", globalPhone, 1)
		if err != nil {
			vcl.ShowMessageFmt("提交错误：%v", err)
			return
		}
		s = "提交成功！\n"
	}
	ProjectDir = ""
	f.StatusBar.SetSimpleText("就绪")
	workload, err := TodayWorkload()
	if err != nil {
		vcl.ShowMessageFmt("%s获取今日已提交任务数失败：%v", s, err)
		return
	}
	MissionInProgress = false
	f.BtnGetMission.SetEnabled(true)

	vcl.ShowMessageFmt("%s今日任务量：\n扫描：%d\n修图：%d\n返工：%d", s, workload.Scan, workload.Edit, workload.Rework)
}
