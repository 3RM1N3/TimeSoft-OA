package lib

import "errors"

type PacketType byte // 定义数据包的类型
type ReportCode byte // 远端的返回码

const (
	SendHead    PacketType = iota // 请求上传文件
	ReceiveHead                   // 请求下载文件
	Report                        // 汇报，用于汇报数据接受情况
	Notice                        // 用于通知
	PullList                      // 客户端请求获得数据表
	PushList                      // 服务端用于推送表
	Login                         // 用于登录
	Signup                        // 用于注册新账号
	ClientCo                      // 用于获取现有的客户公司
)

const (
	Failed          ReportCode = iota // 失败
	Success                           // 成功
	WrongIDOrPwd                      // 用户名或密码错误
	DBOperateErr                      // 数据库操作发生错误
	Unlogin                           // 用户未登录
	ExistingAccount                   // 已存在的账号
)

var codeErrMap = map[ReportCode]error{
	Failed:          errors.New("远端发生错误"),
	WrongIDOrPwd:    errors.New("账号或密码错误"),
	DBOperateErr:    errors.New("数据库操作错误"),
	Unlogin:         errors.New("尚未登录"),
	ExistingAccount: errors.New("账号已存在"),
	Success:         nil,
}

func (c ReportCode) Pack() []byte {
	return []byte{byte(c)}
}

func (c ReportCode) ToError() error {
	e, ok := codeErrMap[c]
	if !ok {
		return errors.New("与远端沟通时发生了不可预料错误")
	}
	return e
}
