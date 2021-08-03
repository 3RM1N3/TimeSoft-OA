package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	globalName, globalID, globalIPAddr string
	conn                               *net.TCPConn
)
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
