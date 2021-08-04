package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var conn *net.TCPConn
var loginSuccess = make(chan bool, 1)

type ScanedJob struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 注册账号
func SignUpAccount(address string, signupJson SP.SignUpJson) error {
	// 创建连接
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("读取地址失败")
		return err
	}
	conn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Println("连接失败")
		return err
	}
	defer conn.Close()

	// 发送数据
	signupData, _ := json.Marshal(signupJson)
	signupData = append([]byte{byte(SP.Signup)}, signupData...)
	n, err := conn.Write(signupData)
	if err != nil || n == 0 {
		fmt.Println("发送数据失败")
		return err
	}

	// 接收数据
	data := make([]byte, 1024)
	n, _, err = conn.ReadFromUDP(data)
	if err != nil || n < 2 {
		fmt.Println("读取数据失败")
		return err
	}

	var report SP.ReportJson
	err = json.Unmarshal(data[:n], &report)
	if err != nil {
		log.Println("读取服务器消息失败")
		return err
	}

	// 处理数据
	if !report.Success {
		log.Println("注册失败")
		return errors.New(report.Msg)
	}

	return nil
}

// 用户登录
func Login(address string, loginJson SP.LoginJson) error {
	// 创建连接
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("读取地址失败", err)
		return err
	}
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Println("连接失败", err)
		return err
	}
	defer udpConn.Close()

	// 发送数据
	loginData, _ := json.Marshal(loginJson)
	loginData = append([]byte{byte(SP.Login)}, loginData...)
	n, err := udpConn.Write(loginData)
	if err != nil || n == 0 {
		fmt.Println("发送数据失败", err)
		return err
	}

	// 接收数据
	data := make([]byte, 1024)
	n, _, err = udpConn.ReadFromUDP(data)
	if err != nil || n < 2 {
		fmt.Println("读取数据失败", err)
		return err
	}

	var report SP.ReportJson
	err = json.Unmarshal(data[:n], &report)
	if err != nil {
		log.Println("读取服务器消息失败", err)
		return err
	}

	// 处理数据
	if !report.Success {
		log.Printf("登录失败：%s\n", report.Msg)
		return errors.New(report.Msg)
	}

	conn, err = ConnectToServer(strings.Replace(address, "8080", "8888", 1)) // 登录至服务器
	if err != nil {
		log.Printf("登录失败")
		return err
	}
	loginSuccess <- true
	close(loginSuccess)
	return nil
}

// 上传文件
func UploadFile(fileName string, sendSPChan chan SP.SocketPacket) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Println("文件打开失败")
		return err
	}
	stat, _ := f.Stat()
	fmt.Printf("准备发送%d字节，", stat.Size())

	fmt.Println("制作为包传入channel")
	if err := SP.NewZipPacket(f, sendSPChan); err != nil {
		log.Println("文件打包失败")
		return err
	}

	return nil
}

// 处理收到的包
func ProcessPacket(receiveSPChan chan SP.SocketPacket) {
	for {
		sp := <-receiveSPChan
		switch sp.TypeByte {
		case SP.FileUpload:
			func() {}()
		}
	}
}

// 连接服务器
func ConnectToServer(address string) (*net.TCPConn, error) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", address)

	conn, err := net.DialTCP("tcp", nil, tcpAddr) // 创建连接
	if err != nil {
		log.Println("拨号失败")
		return nil, err
	}

	conn.SetKeepAlive(true) // 发送心跳包维持连接

	return conn, nil
}

// 设置项目文件夹
func SetProjectDir(dirPath string) error {
	dirEntryList, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("读取目标文件夹失败")
		return err
	}
	if len(dirEntryList) == 0 { // 空文件夹能够直接使用
		return nil
	}

	// 验证.verf文件
	if r, err := CheckVerf(dirPath); err != nil {
		return err
	} else if !r {
		return errors.New("文件夹校验失败，请选择曾用的项目文件夹或空文件夹")
	}

	return nil
}

// 校验.verf文件
func CheckVerf(dirPath string) (bool, error) {
	verfPath := path.Join(dirPath, ".verf")
	f, err := os.Open(verfPath)
	if err != nil {
		log.Println("校验文件打开失败")
		return false, err
	}

	b := make([]byte, 40)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	wantVerfStr := string(b)
	gotVerfStr, err := GenVerf(dirPath)
	if err != nil {
		log.Println("生成校验码失败")
		return false, err
	}
	if gotVerfStr != wantVerfStr {
		return false, nil
	}

	return true, nil
}

// 生成.verf校验字符串
func GenVerf(dirPath string) (string, error) {
	pathList := []string{}

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == ".verf" {
			return nil
		}

		pathList = append(pathList, strings.TrimPrefix(path, dirPath))
		return nil
	})
	if err != nil {
		return "", err
	}

	macs, err := GetMacAddrs()
	if err != nil {
		return "", err
	}
	sort.Strings(macs)

	sort.Strings(pathList)
	got := "!g*657#JW@$" + macs[0] + strings.Join(pathList, "")

	return "verify!#" + SP.MD5(got), nil
}

// 获取本机MAC地址
func GetMacAddrs() ([]string, error) {
	macAddrs := []string{}

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
