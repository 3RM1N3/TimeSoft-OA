package lib

// LoginJson 用户名密码登录json
type LoginJson struct {
	PhoneNumber string
	Pwd         string
}

// SignUpJson 注册账号json
type SignUpJson struct {
	PhoneNumber string
	Pwd         string
	RealName    string
}

// WorkLoadJson 工作量统计
type WorkLoadJson struct {
	Phone  string // 用户账号
	Scan   int
	Edit   int
	Rework int
}
