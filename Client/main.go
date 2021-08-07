package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"TimeSoft-OA/lib"
)

func main() {
	for {
		fmt.Print("\n->")
		s := ""
		fmt.Scanln(&s)
		s = strings.TrimSpace(s)
		switch s {
		case "exit":
			return
		case "send":
			err := lib.SendFile("a.zip", &net.TCPConn{})
			if err != nil {
				log.Println("上传失败", err)
				continue
			}
			fmt.Println("发送成功")

		case "login":
			loginjson := lib.LoginJson{
				PhoneNumber: "13284030601",
				Pwd:         lib.MD5("admin"),
			}
			err := Login("127.0.0.1:8080", loginjson)
			if err != nil {
				log.Println("登录失败", err)
				continue
			}
			fmt.Println("登录成功！")

		case "signup":
			signup := lib.SignUpJson{
				PhoneNumber: "18512341234",
				Pwd:         lib.MD5("admin"),
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
