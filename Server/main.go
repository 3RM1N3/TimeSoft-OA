package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	SP "TimeSoft-OA/SocketPacket"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./timesoft.db")
	if err != nil {
		log.Println("连接数据库失败", err)
		os.Exit(1)
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8888") // 创建用于监听的 socket
	if err != nil {
		log.Println("listen err=", err)
		return
	}
	fmt.Println("开始监听...")

	defer listen.Close() // 服务器结束前关闭 listener

	// 循环等待客户端来链接
	for {
		fmt.Println("阻塞等待客户端连接...")
		conn, err := listen.Accept() // 创建用户数据通信的socket
		if err != nil {
			log.Println("Accept() err=", err)
			continue
		}
		fmt.Printf("与客户端%s连接成功\n", conn.RemoteAddr().String())
		go process(conn) // 起一个协程，为该客户端服务
	}

}

func process(conn net.Conn) {
	defer conn.Close()
	receivedSPChan := make(chan SP.SocketPacket, 128)
	sendSPChan := make(chan SP.SocketPacket, 64)

	go SP.Send(conn, sendSPChan)                       // 发送数据
	go ProcessPacket(conn, sendSPChan, receivedSPChan) // 处理收到的字节数据

	SP.Receive(conn, receivedSPChan)
}

// 处理收到的包
func ProcessPacket(conn net.Conn, sendSPChan, receiveSPChan chan SP.SocketPacket) {
	fileName := ""
	realName := ""
	var fileSize, writeSize int64

	for {
		sp := <-receiveSPChan

		switch sp.TypeByte {
		// 登录
		case SP.Login:
			loginJson := SP.LoginJson{}
			err := json.Unmarshal(sp.Data, &loginJson) // 获取客户端提供的用户名和密码
			if err != nil {
				log.Println("获取用户名或密码失败", err)
				SP.ReportSuccess(false, "获取用户名或密码失败", sendSPChan)
				conn.Close()
				return
			}

			// 从数据库查询对应密码
			row := db.QueryRow(`SELECT PWD, REALNAME FROM USER WHERE USERNAME = ?`, loginJson.User)
			truePwd := ""
			err = row.Scan(&truePwd, &realName)
			if err != nil {
				log.Println("从数据库查询失败", err)
				SP.ReportSuccess(false, "从数据库查询失败", sendSPChan)
				conn.Close()
				return
			}

			// 判断用户名与密码
			if loginJson.Pwd != truePwd {
				SP.ReportSuccess(false, "用户名或密码错误", sendSPChan)
				conn.Close()
				return
			}

			SP.ReportSuccess(true, "", sendSPChan)

		// 客户端上传文件信息至服务器
		case SP.FileUpload:
			if realName == "" {
				log.Println("用户未登录")
				continue
			}
			fileUploadJson := SP.FileUploadJson{}
			json.Unmarshal(sp.Data, &fileUploadJson)
			fileName = fileUploadJson.FileName
			fileSize = fileUploadJson.FileSize
			writeSize = 0

		// 保存文件
		case SP.ZipArchive:
			if realName == "" {
				log.Println("用户未登录")
				continue
			}

			// 在本地创建文件
			fmt.Printf("\n***此次为文件 %s 的第%d部分***\n", fileName, sp.CurrentPart)
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				log.Println("文件打开失败", err)
				continue
			}

			n, err := f.Write(sp.Data) // 写入文件
			writeSize += int64(n)
			if err != nil {
				log.Println("文件写入失败", err)
			}

			f.Close()

			if sp.CurrentPart != sp.AllPart { // 如果文件未结束继续读取
				continue
			}

			fmt.Printf("文件接收完毕\n")

			if fileSize != writeSize { // 写入的总字节数不等于文件实际字节数
				log.Println("接收文件有误，删除该文件")
				os.Remove(fileName)
				fileName = ""
				continue
			}

			fileName = ""
			//err = SP.Unzip(fileName, "./项目文件夹") // 解压文件
			//if err != nil {
			//	log.Println("文件解压失败", err)
			//	os.Remove(fileName)
			//}
		}
	}
}
