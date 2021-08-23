package main

import (
	"TimeSoft-OA/lib"
	"fmt"
	"regexp"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type TSignupForm struct {
	*vcl.TForm
	ButtonLogin    *vcl.TButton
	EditServerAddr *vcl.TEdit
	EditUser       *vcl.TEdit
	Label1         *vcl.TLabel
	EditPwd        *vcl.TEdit
	EditPwdTwice   *vcl.TEdit
	EditRealName   *vcl.TEdit
	Label2         *vcl.TLabel
	Label3         *vcl.TLabel
	Label4         *vcl.TLabel
	Label5         *vcl.TLabel
}

func (f *TSignupForm) OnFormCreate(sender vcl.IObject) {
	// 注册窗口
	f.SetBorderIcons(types.NewSet(types.BiSystemMenu))
	f.SetBorderStyle(types.BsSingle)
	f.SetCaption("注册")
	f.SetPosition(types.PoDesktopCenter)
	f.SetHeight(385)
	f.SetWidth(300)

	// 标签：“服务器：”
	f.Label1 = vcl.NewLabel(f)
	f.Label1.SetParent(f)
	f.Label1.SetCaption("服务器：")
	label1Font := vcl.NewFont()  // 初始化字体
	label1Font.SetSize(10)       // 设置字体大小
	f.Label1.SetFont(label1Font) // 设置字体
	f.Label1.SetHeight(20)
	f.Label1.SetLeft(10)
	f.Label1.SetTop(10)
	f.Label1.SetWidth(52)

	// 输入框：服务器地址
	f.EditServerAddr = vcl.NewEdit(f)
	f.EditServerAddr.SetParent(f)
	editFont := vcl.NewFont()          // 初始化字体
	editFont.SetSize(12)               // 设置字体大小
	f.EditServerAddr.SetFont(editFont) // 设置字体
	f.EditServerAddr.SetMaxLength(15)
	f.EditServerAddr.SetTextHint("0.0.0.0")
	f.EditServerAddr.SetHeight(30)
	f.EditServerAddr.SetLeft(10)
	f.EditServerAddr.SetTop(35)
	f.EditServerAddr.SetWidth(280)

	// 标签：“账　号：”
	f.Label2 = vcl.NewLabel(f)
	f.Label2.SetParent(f)
	f.Label2.SetCaption("填写电话号码作为账号：")
	f.Label2.SetFont(label1Font) // 设置字体
	f.Label2.SetHeight(20)
	f.Label2.SetLeft(10)
	f.Label2.SetTop(75)
	f.Label2.SetWidth(52)

	// 输入框：账号
	f.EditUser = vcl.NewEdit(f)
	f.EditUser.SetParent(f)
	f.EditUser.SetFont(editFont) // 设置字体
	f.EditUser.SetMaxLength(11)
	f.EditUser.SetHeight(30)
	f.EditUser.SetLeft(10)
	f.EditUser.SetTop(100)
	f.EditUser.SetWidth(280)

	// 标签：“请输入姓名：”
	f.Label5 = vcl.NewLabel(f)
	f.Label5.SetParent(f)
	f.Label5.SetCaption("请输入姓名：")
	f.Label5.SetFont(label1Font) // 设置字体
	f.Label5.SetHeight(20)
	f.Label5.SetLeft(10)
	f.Label5.SetTop(140)
	f.Label5.SetWidth(52)

	// 输入框：姓名
	f.EditRealName = vcl.NewEdit(f)
	f.EditRealName.SetParent(f)
	f.EditRealName.SetFont(editFont) // 设置字体
	f.EditRealName.SetMaxLength(11)
	f.EditRealName.SetHeight(30)
	f.EditRealName.SetLeft(10)
	f.EditRealName.SetTop(165)
	f.EditRealName.SetWidth(280)

	// 标签：“密　码：”
	f.Label3 = vcl.NewLabel(f)
	f.Label3.SetParent(f)
	f.Label3.SetCaption("请输入密码：")
	f.Label3.SetFont(label1Font) // 设置字体
	f.Label3.SetHeight(20)
	f.Label3.SetLeft(10)
	f.Label3.SetTop(205)
	f.Label3.SetWidth(52)

	// 输入框：密码
	f.EditPwd = vcl.NewEdit(f)
	f.EditPwd.SetParent(f)
	f.EditPwd.SetFont(editFont) // 设置字体
	f.EditPwd.SetPasswordChar('*')
	f.EditPwd.SetMaxLength(18)
	f.EditPwd.SetHeight(30)
	f.EditPwd.SetLeft(10)
	f.EditPwd.SetTop(230)
	f.EditPwd.SetWidth(280)

	// 标签：“验　证：”
	f.Label4 = vcl.NewLabel(f)
	f.Label4.SetParent(f)
	f.Label4.SetCaption("请再次输入密码：")
	f.Label4.SetFont(label1Font) // 设置字体
	f.Label4.SetHeight(20)
	f.Label4.SetLeft(10)
	f.Label4.SetTop(270)
	f.Label4.SetWidth(52)

	// 输入框：验证密码
	f.EditPwdTwice = vcl.NewEdit(f)
	f.EditPwdTwice.SetParent(f)
	f.EditPwdTwice.SetFont(editFont) // 设置字体
	f.EditPwdTwice.SetPasswordChar('*')
	f.EditPwdTwice.SetMaxLength(18)
	f.EditPwdTwice.SetHeight(30)
	f.EditPwdTwice.SetLeft(10)
	f.EditPwdTwice.SetTop(295)
	f.EditPwdTwice.SetWidth(280)

	// 提交按钮
	f.ButtonLogin = vcl.NewButton(f)
	f.ButtonLogin.SetParent(f)
	f.ButtonLogin.SetCaption("提交")
	loginBtnFont := vcl.NewFont()       // 初始化字体
	loginBtnFont.SetSize(11)            // 设置字体大小
	f.ButtonLogin.SetFont(loginBtnFont) // 设置字体
	f.ButtonLogin.SetHeight(40)
	f.ButtonLogin.SetLeft(10)
	f.ButtonLogin.SetTop(335)
	f.ButtonLogin.SetWidth(280)
	f.ButtonLogin.SetOnClick(f.OnSubmitClick)
}

// 点击提交按钮
func (f *TSignupForm) OnSubmitClick(sender vcl.IObject) {
	// 验证服务器地址合法性
	serverAddr := f.EditServerAddr.Text()
	if serverAddr == "" {
		vcl.ShowMessage("服务器地址不能为空")
		return
	}
	re := regexp.MustCompile(`\d{1,3}(\.\d{1,3}){3}`)
	if serverAddr != re.FindString(serverAddr) {
		vcl.ShowMessage("服务器地址格式不合法")
		return
	}
	globalServerAddr = serverAddr

	// 验证电话号码合法性
	phone := f.EditUser.Text()
	if phone == "" {
		vcl.ShowMessage("账号不能为空")
		return
	}
	re = regexp.MustCompile(`1\d{10}`)
	if phone != re.FindString(phone) {
		vcl.ShowMessage("电话号码格式不合法")
		return
	}

	// 验证姓名合法性
	name := f.EditRealName.Text()
	if name == "" {
		vcl.ShowMessage("姓名不能为空")
		return
	}
	re = regexp.MustCompile("[\u4e00-\u9fa5]{2,5}")
	if name != re.FindString(name) {
		vcl.ShowMessage("姓名格式不合法，应为2-5个汉字")
		return
	}

	pwd := f.EditPwd.Text()
	if pwd == "" {
		vcl.ShowMessage("密码不能为空")
		return
	}

	twicePwd := f.EditPwdTwice.Text()
	if twicePwd != pwd {
		vcl.ShowMessage("两次输入的密码不一致")
		return
	}

	signupjson := lib.SignUpJson{
		PhoneNumber: globalPhone,
		Pwd:         lib.MD5(pwd),
		RealName:    name,
	}
	err := SignUpAccount(globalServerAddr, signupjson)
	if err != nil {
		vcl.ShowMessageFmt("注册失败：%v", err)
		return
	}
	fmt.Println("注册成功！")
	f.Close()
}

func (f *TSignupForm) OnFormCloseQuery(Sender vcl.IObject, CanClose *bool) {
	LoginForm.Show()
	*CanClose = true
}
