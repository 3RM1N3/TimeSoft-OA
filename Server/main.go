package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./timesoft.db")
	if err != nil {
		log.Println("连接数据库失败", err)
		os.Exit(1)
	}
}

func main() {
	defer db.Close()

	go SignupAndLogin() // 开启用于注册和登录的udp服务器

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
		//fmt.Printf("与客户端%s连接成功\n", conn.RemoteAddr().String())
		go process(conn) // 起一个协程，为该客户端服务
	}
}
