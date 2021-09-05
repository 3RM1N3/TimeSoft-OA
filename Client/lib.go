package main

import (
	"TimeSoft-OA/lib"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
)

type ScanedJob struct {
	JobID        string `json:"jobid"` // 任务ID
	FolderName   string `json:"foldername"` // 文件夹名
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 全局变量

var ( // config.json 控制
	globalTCPPort        = ":8888" // 远程TCP服务器端口
	globalUDPPort        = ":8080" // 远程UDP服务器端口
	globalReworkInterval = 1       // 返工任务刷新时间间隔，单位小时
)

var ( // 错误信息
	ErrScanTooFast = errors.New("你的扫描速度似乎有些快于常人，为判定是否作弊，请主动与管理员取得联系。")
	ErrFindNotJpg  = errors.New("检测到非*.jpg格式文件，请确认扫描设置正确，删除格式错误的文件后重试。")
)

var ( // 实例化窗体
	SignupForm *TSignupForm // 注册窗体
	LoginForm  *TLoginForm  // 登录窗体
	MainForm   *TMainForm   // 主窗体
)

var globalPhone string       // 全局变量用户账号
var globalServerAddr string  // 全局服务器地址
var macAddr = ""             // 本机mac地址
var ProjectDir = ""          // 全局项目文件夹
var MissionList []string     // 任务列表

var MissionInProgress = false        // 是否有任务进行中
var OverScan = false                 // 扫描结束
var ChStartScan = make(chan bool, 2) // 开始监测项目文件夹

// SignUpAccount 注册账号，传入服务器地址和注册结构体，返回错误
func SignUpAccount(address string, signupJson lib.SignUpJson) error {
	b, err := sendUDPMsg(address, lib.Signup, signupJson)
	if err != nil {
		return err
	}

	return lib.ReportCode(b[0]).ToError()
}

// Login 用户登录，传入服务器地址和登录结构体，返回错误
func Login(address string, loginJson lib.LoginJson) error {
	b, err := sendUDPMsg(address, lib.Login, loginJson)
	if err != nil {
		return err
	}

	err = lib.ReportCode(b[0]).ToError()
	if err == nil {
		globalPhone = loginJson.PhoneNumber
	}

	return err
}

// 从服务器获取管理员分配的修图任务，返回任务数量和错误
func getEditMission() (uint32, error) {
	var editMissions = []string{globalPhone}
	b, err := sendUDPMsg(globalServerAddr, lib.EditMission, editMissions)
	if err != nil {
		return 0, err
	}

	if len(b) == 1 {
		return 0, lib.ReportCode(b[0]).ToError()
	}

	return lib.ByteToUint32(b)
}

// 发送udp消息，jsonStruct为要发送的结构体，返回收到的字节切片和错误类型；
// address 为不包含端口的服务器地址，端口号从全局变量 globalUDPPort 获取；
// packetType为常量包类型
func sendUDPMsg(address string, packType lib.PacketType, jsonStruct interface{}) ([]byte, error) {
	// 创建连接
	udpAddr, err := net.ResolveUDPAddr("udp", address+globalUDPPort)
	if err != nil {
		log.Println("读取地址失败")
		return nil, err
	}
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Println("连接失败")
		return nil, err
	}
	defer udpConn.Close()

	// 识别类型并发送数据
	var jsonData []byte
	switch jsonStruct.(type) {
	case nil:
		break
	default:
		jsonData, _ = json.Marshal(jsonStruct)
	}
	jsonData = append([]byte{byte(packType)}, jsonData...)
	n, err := udpConn.Write(jsonData)
	if err != nil || n == 0 {
		fmt.Println("发送数据失败")
		return nil, err
	}

	// 接收返回消息
	reportBytes := make([]byte, 16384)
	n, _, err = udpConn.ReadFromUDP(reportBytes)
	if err != nil {
		fmt.Println("读取数据失败")
		return nil, err
	}

	return reportBytes[:n], nil
}

// GetMacAddrs 获取本机MAC地址，返回全部mac地址的字符串列表和错误
func GetMacAddrs() ([]string, error) {
	var macAddrs []string

	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if macAddr != "" {
			macAddrs = append(macAddrs, macAddr)
		}
	}
	return macAddrs, nil
}

// ClientReceiveFile 客户端从远端接收文件，第一个参数弃用，传入一个net.Conn类型；
// 返回文件名和错误
func ClientReceiveFile(fileList []string, conn net.Conn) (string, error) {
	defer conn.Close()

	downloadHead := lib.FileReceiveHead{
		FileList:   fileList,
		Downloader: globalPhone,
	}
	fileHead, err := downloadHead.MakeHead()
	if err != nil {
		return "", err
	}

	_, err = conn.Write(fileHead)
	if err != nil {
		return "", err
	}
	fmt.Println("发送文件头成功")

	result := make([]byte, 1)
	conn.Read(result)
	if result[0] != '1' {
		return "", errors.New("服务器端发送文件时发生错误")
	}
	fmt.Println("可以接收正文")

	// 接收正文
	conn.Read(make([]byte, 1)) // 抛弃即将接收的包头类型（1字节）
	sendHead, err := lib.ReceiveFile(conn)
	if err != nil {
		return "", err
	}
	return sendHead.Name, nil
}

// ReworkItem 修图者设置返工任务，传入档号，如有错误返回错误类型
func ReworkItem(fileID string) error {
	b, err := sendUDPMsg(globalServerAddr, lib.ReworkItem, []string{fileID})
	if err != nil {
		return err
	}

	return lib.ReportCode(b[0]).ToError()
}

// VerifyStringRe 验证字符串是否完全符合正则表达式
// 传入正则字符串和待检测字符串，返回bool值
func VerifyStringRe(reString, dstString string) bool {
	if dstString == "" {
		return false
	}
	re := regexp.MustCompile(reString)
	return dstString == re.FindString(dstString)
}