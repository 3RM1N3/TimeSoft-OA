package main

import (
	"TimeSoft-OA/lib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var globalPhone string      // 全局变量用户账号
var globalServerAddr string // 全局服务器地址
var (
	globalTCPPort        = ":8888" // 远程TCP服务器端口
	globalUDPPort        = ":8080" // 远程UDP服务器端口
	globalReworkInterval = 1       // 返工任务刷新时间间隔，单位小时
)

type ScanedJob struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 读取配置文件config.json
func init() {
	// 定义全局变量结构体
	config := struct {
		GlobalTCPPort        string `json:"远程TCP服务器端口"`
		GlobalUDPPort        string `json:"远程UDP服务器端口"`
		GlobalReworkInterval int    `json:"返工任务刷新时间间隔/小时"`
	}{
		GlobalTCPPort:        globalTCPPort,
		GlobalUDPPort:        globalUDPPort,
		GlobalReworkInterval: globalReworkInterval,
	}

	// 检查配置文件存在与否
	configFile := "config.json"
	stat, err := os.Stat(configFile)

	if os.IsNotExist(err) { // 配置文件不存在，创建配置文件
		f, err := os.Create(configFile)
		if err != nil {
			log.Printf("创建配置文件失败：%v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		b, _ := json.MarshalIndent(&config, "", "    ")
		_, err = f.Write(b)
		if err != nil {
			log.Println("写入配置文件失败")
			os.Exit(1)
		}
		return
	}

	// 配置文件存在
	if stat.IsDir() { // 配置文件名被占用
		log.Println("错误：文件名“config.json”被文件夹占用，请删除或重命名该文件夹后重试")
		os.Exit(1)
	}

	f, err := os.Open(configFile) // 打开配置文件
	if err != nil {
		log.Println("读取配置文件失败，使用默认参数")
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f) // 读取配置文件
	if err != nil {
		log.Println("读取配置文件失败，使用默认参数")
		return
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Println("配置文件包含语法错误，使用默认参数")
		return
	}

	// 设置全局变量
	globalTCPPort = config.GlobalTCPPort
	globalUDPPort = config.GlobalUDPPort
	globalReworkInterval = config.GlobalReworkInterval
}

// 注册账号
func SignUpAccount(address string, signupJson lib.SignUpJson) error {
	b, err := sendUDPMsg(address, lib.Signup, signupJson)
	if err != nil {
		return err
	}

	return lib.ReportCode(b[0]).ToError()
}

// 用户登录
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

// 发送udp消息，jsonStruct为要发送的结构体，返回收到的字节切片和错误类型；
// address 为不包含端口的服务器地址，端口号从全局变量 globalUDPPort 获取
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
	jsonData := []byte{}
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
	reportBytes := make([]byte, 1024)
	n, _, err = udpConn.ReadFromUDP(reportBytes)
	if err != nil {
		fmt.Println("读取数据失败")
		return nil, err
	}

	return reportBytes[:n], nil
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

// 客户端从远端接收文件
func ClientReceiveFile(fileList []string, conn *net.TCPConn) (string, error) {
	defer conn.Close()

	downloadHead := lib.FileReceiveHead{
		FileList:   fileList,
		Downloader: "13284030601",
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
