package main

import (
	"database/sql"
	"log"
	"os"
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./timesoft.db") // 连接数据库
	if err != nil {
		log.Println("连接数据库失败", err)
		os.Exit(1)
	}
}

func main() {
	defer db.Close()

	go UDPServer(PortUDP) // 开启用于注册和登录的UDP服务器

	err := TCPServer(PortTCP) // 开启收发文件的TCP服务器
	if err != nil {
		log.Println("TCP服务器开启失败", err)
	}
}
