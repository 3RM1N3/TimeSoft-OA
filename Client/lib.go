package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	globalName, globalID, globalIPAddr string
)

type ScanedJob struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 用户登录
func Login(conn net.Conn, userName, pwd string) error {
	pwdMD5 := SP.MD5(pwd)

	// 发送登录信息
	loginJson := SP.LoginJson{
		User: userName,
		Pwd:  pwdMD5,
	}
	loginJsonData, _ := json.Marshal(loginJson)
	SP.NewJsonPacket(SP.Login, loginJsonData, sendSPChan)

	// 获取返回信息
	sp := <-receiveSPChan
	if sp.TypeByte != SP.Report {
		return errors.New("返回信息错误")
	}

	reportJson := SP.ReportJson{}
	err := json.Unmarshal(sp.Data, &reportJson)
	if err != nil {
		return err
	}

	if !reportJson.Success {
		return errors.New(reportJson.Msg)
	}

	return nil
}

func UploadFile(fileName string, sendSPChan chan SP.SocketPacket) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	stat, _ := f.Stat()
	fmt.Printf("准备发送%d字节，", stat.Size())

	fmt.Println("制作为包传入channel")
	if err := SP.NewZipPacket(f, sendSPChan); err != nil {
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
		log.Println("拨号失败", err)
		return nil, err
	}

	conn.SetKeepAlive(true) // 发送心跳包维持连接

	return conn, nil
}
