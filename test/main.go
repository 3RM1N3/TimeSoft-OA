package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"
)

func main() {
	//a()
	//b()
	SignupAndLogin()
}

// 用于注册和登录的udp服务器
func SignupAndLogin() {
	//建立一个UDP监听
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8080,
	})
	if err != nil {
		log.Printf("登录和注册服务未开启, err:%v\n", err)
		os.Exit(1)
		return
	}
	log.Println("登录和注册服务已开启")
	defer listen.Close()

	// 使用for循环监听消息
	for {
		bytePack := make([]byte, 1024) // 初始化保存接收数据的变量
		reportSP := SP.ReportJson{     // 初始化一个汇报json
			Success: false,
			Msg:     "",
		}

		// 读取UDP数据
		n, addr, err := listen.ReadFromUDP(bytePack)
		if err != nil || n < 2 {
			log.Printf("read udp failed: %v\n", err)
			continue
		}

		// 处理数据
		if bytePack[0] == 0x0 { // 注册
			var signup SP.SignUpJson
			err := json.Unmarshal(bytePack[1:n], &signup)
			if err != nil {
				reportSP.Success = false
				reportSP.Msg = "发送的数据有误"
			} else {
				log.Printf("%s 申请注册\n", addr)
				fmt.Printf("    密码：%s\n    真实姓名：%s\n    电话号：%s\n", signup.Pwd, signup.RealName, signup.PhoneNumber)
				reportSP.Success = true
				reportSP.Msg = ""
			}

		} else if bytePack[0] == 0x1 { // 登录
			reportSP.Success = true
			reportSP.Msg = ""
		} else { // 消息有误
			reportSP.Success = false
			reportSP.Msg = "发送的数据有误"
		}

		// 返回数据
		reportBytes, _ := json.Marshal(reportSP)
		_, err = listen.WriteToUDP(reportBytes, addr)
		if err != nil {
			log.Printf("write udp failed, err:%v\n", err)
			continue
		}
	}
}

// 测试解压文件
func b() {
	fmt.Println(path.Dir("a.txt"))
	err := SP.Unzip("a48.6.zip", "./folder")
	if err != nil {
		log.Println("解压失败", err)
	}
}
