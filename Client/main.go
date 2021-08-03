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
	go func() {
		for {
			if <-loginSuccess {
				break
			}
		}

		fmt.Printf("登录成功！\n")

		go SP.Receive(conn, receiveSPChan) // 接收TCP数据
		go SP.Send(conn, sendSPChan)       // 发送TCP数据
		ProcessPacket(receiveSPChan)       // 处理收到的TCP数据
	}()

	for {
		fmt.Print("\n->")
		s := ""
		fmt.Scanln(&s)
		s = strings.TrimSpace(s)
		switch s {
		case "exit":
			return
		case "send":
			err := UploadFile("a.zip", sendSPChan)
			if err != nil {
				log.Println("上传失败", err)
				continue
			}
			fmt.Println("发送成功")

		case "login":
			loginjson := SP.LoginJson{
				PhoneNumber: "13284030601",
				Pwd:         SP.MD5("admin"),
			}
			err := Login("127.0.0.1:8080", loginjson)
			if err != nil {
				log.Println("登录失败", err)
				continue
			}
			fmt.Println("登录成功！")

		case "signup":
			signup := SP.SignUpJson{
				PhoneNumber: "18512341234",
				Pwd:         SP.MD5("admin"),
				RealName:    "李雷",
			}
			err := SignUpAccount("127.0.0.1:8080", signup)
			if err != nil {
				log.Println("注册失败", err)
				continue
			}
			fmt.Println("注册成功！")
		}
	}
}
