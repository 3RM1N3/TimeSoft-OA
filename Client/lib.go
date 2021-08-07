package main

import (
	"TimeSoft-OA/lib"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

var globalPhone string      // 全局变量用户账号
var globalServerAddr string // 全局服务器地址

type ScanedJob struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 注册账号
func SignUpAccount(address string, signupJson lib.SignUpJson) error {
	return sendUDPMsg(address, lib.Signup, signupJson)
}

// 用户登录
func Login(address string, loginJson lib.LoginJson) error {
	err := sendUDPMsg(address, lib.Login, loginJson)
	if err != nil {
		return err
	}

	globalPhone = loginJson.PhoneNumber
	return nil
}

// 发送udp消息，jsonStruct为要发送的结构体
func sendUDPMsg(address string, packType lib.PacketType, jsonStruct interface{}) error {
	// 创建连接
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("读取地址失败")
		return err
	}
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Println("连接失败")
		return err
	}
	defer udpConn.Close()

	// 发送数据
	jsonData, _ := json.Marshal(jsonStruct)
	jsonData = append([]byte{byte(packType)}, jsonData...)
	n, err := udpConn.Write(jsonData)
	if err != nil || n == 0 {
		fmt.Println("发送数据失败")
		return err
	}

	// 接收返回消息
	reportByte := make([]byte, 1)
	_, _, err = udpConn.ReadFromUDP(reportByte)
	if err != nil {
		fmt.Println("读取数据失败")
		return err
	}

	return lib.ReportCode(reportByte[0]).ToError()
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
func ClientReceiveFile(fileList []string, conn *net.TCPConn) (string, *os.File, error) {
	defer conn.Close()

	downloadHead := lib.FileReceiveHead{
		FileList:   fileList,
		Downloader: "13284030601",
	}
	fileHead, err := downloadHead.MakeHead()
	if err != nil {
		return "", nil, err
	}

	_, err = conn.Write(fileHead)
	if err != nil {
		return "", nil, err
	}
	fmt.Println("发送文件头成功")

	result := make([]byte, 1)
	conn.Read(result)
	if result[0] != '1' {
		return "", nil, errors.New("服务器端发送文件时发生错误")
	}
	fmt.Println("可以接收正文")

	// 接收正文
	conn.Read(make([]byte, 1)) // 抛弃即将接收的包头类型（1字节）
	return lib.ReceiveFile(conn)
}
