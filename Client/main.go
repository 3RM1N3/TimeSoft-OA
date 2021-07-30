package main

import (
	"fmt"
	"log"
	"strings"

	SP "TimeSoft-OA/SocketPacket"
)

var sendSPChan = make(chan SP.SocketPacket, 64)
var receiveSPChan = make(chan SP.SocketPacket, 128)

func main() {
	log.Println("连接服务器...")

	conn, err := ConnectToServer("127.0.0.1:8888")
	if err != nil {
		log.Println("连接失败", err)
		return
	}
	log.Println("连接成功")

	go SP.Receive(conn, receiveSPChan)
	go SP.Send(conn, sendSPChan)

	log.Println("登录中...")
	err = Login(conn, "admin", "admin")
	if err != nil {
		fmt.Println("登录失败：", err)
		return
	}
	log.Println("登录成功")

	// 处理读取到的数据
	go ProcessPacket(receiveSPChan)

	for {
		fmt.Print("\n->")
		s := ""
		fmt.Scanln(&s)
		s = strings.TrimSpace(s)
		switch s {
		case "exit":
			return
		case "send":
			UploadFile("a.zip", sendSPChan)
		}
	}
}
