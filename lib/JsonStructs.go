package lib

// 报告json
type ReportJson struct {
	Success bool
	Msg     string
}

// 用户名密码登录json
type LoginJson struct {
	PhoneNumber string
	Pwd         string
}

// 注册账号json
type SignUpJson struct {
	PhoneNumber string
	Pwd         string
	RealName    string
}
