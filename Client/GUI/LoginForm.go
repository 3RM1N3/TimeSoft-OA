package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

type TLoginForm struct {
	*vcl.TForm
	ButtonLogin    *vcl.TButton
	EditServerAddr *vcl.TEdit
	EditUser       *vcl.TEdit
	Label1         *vcl.TLabel
	EditPwd        *vcl.TEdit
	Label2         *vcl.TLabel `events:"OnLabel1Click"`
	Label3         *vcl.TLabel `events:"OnLabel1Click"`
	LabelTitle     *vcl.TLabel
	LabelSignUp    *vcl.TLabel
	LabelForgetPwd *vcl.TLabel
}

var LoginForm *TLoginForm

func (f *TLoginForm) OnFormCreate(sender vcl.IObject) {
	// 登录窗口
	f.SetBorderIcons(types.NewSet(types.BiSystemMenu, types.BiMinimize))
	f.SetBorderStyle(types.BsSingle)
	f.SetCaption("登录")
	f.SetPosition(types.PoDesktopCenter)
	f.SetHeight(300)
	f.SetWidth(430)

	// 标题文本
	f.LabelTitle = vcl.NewLabel(f)
	f.LabelTitle.SetParent(f)
	f.LabelTitle.SetAlignment(types.TaCenter)
	f.LabelTitle.SetCaption("时源办公自动化管理系统")
	titleFont := vcl.NewFont()      // 初始化字体
	titleFont.SetSize(14)           // 设置字体大小
	f.LabelTitle.SetFont(titleFont) // 设置字体
	f.LabelTitle.SetHeight(25)
	f.LabelTitle.SetLeft(110)
	f.LabelTitle.SetTop(40)
	f.LabelTitle.SetWidth(210)

	// 标签：“服务器：”
	f.Label1 = vcl.NewLabel(f)
	f.Label1.SetParent(f)
	f.Label1.SetCaption("服务器：")
	label1Font := vcl.NewFont()  // 初始化字体
	label1Font.SetSize(10)       // 设置字体大小
	f.Label1.SetFont(label1Font) // 设置字体
	f.Label1.SetHeight(19)
	f.Label1.SetLeft(80)
	f.Label1.SetTop(115)
	f.Label1.SetWidth(52)

	// 输入框：服务器地址
	f.EditServerAddr = vcl.NewEdit(f)
	f.EditServerAddr.SetParent(f)
	addrFont := vcl.NewFont()          // 初始化字体
	addrFont.SetSize(12)               // 设置字体大小
	f.EditServerAddr.SetFont(addrFont) // 设置字体
	f.EditServerAddr.SetMaxLength(20)
	f.EditServerAddr.SetTextHint("0.0.0.0")
	f.EditServerAddr.SetHeight(29)
	f.EditServerAddr.SetLeft(144)
	f.EditServerAddr.SetTop(110)
	f.EditServerAddr.SetWidth(200)

	// 标签：“账　号：”
	f.Label2 = vcl.NewLabel(f)
	f.Label2.SetParent(f)
	f.Label2.SetCaption("账　号：")
	label2Font := vcl.NewFont()  // 初始化字体
	label2Font.SetSize(10)       // 设置字体大小
	f.Label2.SetFont(label2Font) // 设置字体
	f.Label2.SetHeight(19)
	f.Label2.SetLeft(80)
	f.Label2.SetTop(150)
	f.Label2.SetWidth(52)

	// 输入框：账号
	f.EditUser = vcl.NewEdit(f)
	f.EditUser.SetParent(f)
	userFont := vcl.NewFont()    // 初始化字体
	userFont.SetSize(12)         // 设置字体大小
	f.EditUser.SetFont(userFont) // 设置字体
	f.EditUser.SetMaxLength(15)
	f.EditUser.SetHeight(29)
	f.EditUser.SetLeft(144)
	f.EditUser.SetTop(145)
	f.EditUser.SetWidth(200)

	// 标签：“密　码：”
	f.Label3 = vcl.NewLabel(f)
	f.Label3.SetParent(f)
	f.Label3.SetCaption("密　码：")
	label3Font := vcl.NewFont()  // 初始化字体
	label3Font.SetSize(10)       // 设置字体大小
	f.Label3.SetFont(label3Font) // 设置字体
	f.Label3.SetHeight(19)
	f.Label3.SetLeft(80)
	f.Label3.SetTop(185)
	f.Label3.SetWidth(52)

	// 输入框：密码
	f.EditPwd = vcl.NewEdit(f)
	f.EditPwd.SetParent(f)
	pwdFont := vcl.NewFont()   // 初始化字体
	pwdFont.SetSize(12)        // 设置字体大小
	f.EditPwd.SetFont(pwdFont) // 设置字体
	f.EditPwd.SetPasswordChar('*')
	f.EditPwd.SetMaxLength(18)
	f.EditPwd.SetHeight(29)
	f.EditPwd.SetLeft(144)
	f.EditPwd.SetTop(180)
	f.EditPwd.SetWidth(200)

	// 标签：注册账号
	f.LabelSignUp = vcl.NewLabel(f)
	f.LabelSignUp.SetParent(f)
	f.LabelSignUp.SetCaption("注册账号")
	labelSignupFont := vcl.NewFont()        // 初始化字体
	labelSignupFont.SetColor(colors.ClGrey) // 设置字体大小
	f.LabelSignUp.SetFont(labelSignupFont)  // 设置字体
	f.LabelSignUp.SetHeight(17)
	f.LabelSignUp.SetLeft(10)
	f.LabelSignUp.SetTop(273)
	f.LabelSignUp.SetWidth(48)
	f.LabelSignUp.SetOnClick(f.OnSignUpClick)

	// 标签：忘记密码
	f.LabelForgetPwd = vcl.NewLabel(f)
	f.LabelForgetPwd.SetParent(f)
	f.LabelForgetPwd.SetCaption("忘记密码")
	forgetPwdFont := vcl.NewFont()          // 初始化字体
	forgetPwdFont.SetColor(colors.ClGrey)   // 设置字体大小
	f.LabelForgetPwd.SetFont(forgetPwdFont) // 设置字体
	f.LabelForgetPwd.SetHeight(17)
	f.LabelForgetPwd.SetLeft(372)
	f.LabelForgetPwd.SetTop(273)
	f.LabelForgetPwd.SetWidth(48)
	f.LabelForgetPwd.SetOnClick(f.OnForgetPwdClick)

	// 登录按钮
	f.ButtonLogin = vcl.NewButton(f)
	f.ButtonLogin.SetParent(f)
	f.ButtonLogin.SetCaption("登录")
	loginBtnFont := vcl.NewFont()       // 初始化字体
	loginBtnFont.SetSize(11)            // 设置字体大小
	f.ButtonLogin.SetFont(loginBtnFont) // 设置字体
	f.ButtonLogin.SetHeight(40)
	f.ButtonLogin.SetLeft(95)
	f.ButtonLogin.SetTop(240)
	f.ButtonLogin.SetWidth(240)
	f.ButtonLogin.SetAlign(types.AlNone)
	f.ButtonLogin.SetOnClick(f.OnLoginClick)
}
