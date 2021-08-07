package main

import (
	"TimeSoft-OA/lib"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

// 收发文件的TCP服务器
func TCPServer(address string) error {
	listen, err := net.Listen("tcp", address) // 开启tcp服务器
	if err != nil {
		return err
	}

	log.Println("TCP服务器已开启")
	defer listen.Close() // 服务器结束前关闭监听

	// 循环阻塞等待客户端连接
	for {
		conn, err := listen.Accept() // 创建连接
		if err != nil {
			log.Println("建立连接失败:", err)
			continue
		}

		log.Printf("与客户端%s连接成功\n", conn.RemoteAddr().String())
		go processTCPConn(conn) // 后台为该客户端服务
	}
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

	log.Println("UDP服务已开启")
	defer listen.Close() // 服务器结束前关闭监听

	// 使用for循环监听消息
	for {
		udpMessage := make([]byte, 1024) // 初始化保存接收数据的变量

		// 读取UDP数据
		n, addr, err := listen.ReadFromUDP(udpMessage)
		if err != nil {
			log.Printf("read udp failed: %v\n", err)
			continue
		}

		log.Printf("收到%s的消息\n", addr.String())
		go processUDPMsg(udpMessage[:n], addr, listen) // 后台处理此条消息

	}
}

// 处理TCP连接
func processTCPConn(conn net.Conn) {
	defer conn.Close()

	// 读取一个字节判断连接类型
	b := make([]byte, 1)
	_, err := conn.Read(b)
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == byte(lib.ReceiveHead) {
		fName, f, err := lib.ReceiveFile(conn)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("收到文件", fName)

		// 处理文件
		f.Close()

	} else if b[0] == byte(lib.SendHead) {
		ServerSendFile(conn)
	}
	// 否则即远端返回了无法识别的消息，函数也返回
}

// 处理UDP消息
func processUDPMsg(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	switch udpMessage[0] {
	case byte(lib.Signup): // 注册
		processSignin(udpMessage, addr, listen)

	case byte(lib.Login): // 登录
		processLogin(udpMessage, addr, listen)

	default:
		// 消息有误不做处理
		log.Printf("收到 %s 未识别的消息\n", addr.String())
	}
}

// 处理用户注册
func processSignin(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	var signup lib.SignUpJson
	err := json.Unmarshal(udpMessage[1:], &signup)
	if err != nil {
		listen.WriteToUDP(lib.Failed.Pack(), addr) // 发送错误消息
		return
	}
	log.Printf("%s 申请注册\n", addr.String())

	// 判断与数据库中的项目有无重复
	row := db.QueryRow(`SELECT PHONE FROM USER WHERE PHONE = ? UNION SELECT PHONE FROM UNREVIEWED_USER WHERE PHONE = ?`, signup.PhoneNumber, signup.PhoneNumber)
	err = row.Scan()

	if err != sql.ErrNoRows { // 如果有数据，账号已存在
		listen.WriteToUDP(lib.ExistingAccount.Pack(), addr) // 发送错误消息
		return
	}

	// 在已有用户中唯一，可以加入未审核表中
	_, err = db.Exec(`INSERT INTO UNREVIEWED_USER VALUES (?, ?, ?)`, signup.PhoneNumber, signup.Pwd, signup.RealName)
	if err != nil {
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	listen.WriteToUDP(lib.Success.Pack(), addr) // 成功
	log.Printf("新员工注册待审核：\n    电话：%s\n    姓名：%s\n", signup.PhoneNumber, signup.RealName)
}

// 处理用户登录
func processLogin(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	userRealName := ""
	loginJson := lib.LoginJson{}
	err := json.Unmarshal(udpMessage[1:], &loginJson) // 获取客户端提供的用户名和密码
	if err != nil {
		listen.WriteToUDP(lib.Failed.Pack(), addr) // 发送错误消息
		return
	}

	// 从数据库查询对应密码
	row := db.QueryRow(`SELECT PWD, REALNAME FROM USER WHERE PHONE = ?`, loginJson.PhoneNumber)
	truePwd := ""
	err = row.Scan(&truePwd, &userRealName)
	if err != nil {
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	if loginJson.Pwd != truePwd { // 判断用户名与密码
		listen.WriteToUDP(lib.WrongIDOrPwd.Pack(), addr) // 发送错误消息
		return
	}

	// 密码正确
	log.Printf("%s 用户%s登录一次\n", addr.String(), userRealName)
	listen.WriteToUDP(lib.Success.Pack(), addr) // 成功
}

// TCP服务器发送文件给远端
func ServerSendFile(conn net.Conn) error {
	var head lib.FileReceiveHead
	headSize := uint16(0)
	buf := []byte{}

	for {
		b := make([]byte, 524288)
		n, err := conn.Read(b)
		if err != nil {
			return err
		}

		buf = append(buf, b[:n]...)
		if len(buf) < 3 {
			continue
		}

		if headSize == 0 {
			headSize, err = lib.ByteToUint16(buf[:2])
			if err != nil {
				return err
			}
			buf = buf[2:]
		}

		if len(buf) < int(headSize) {
			continue
		}

		// 解析头部
		fmt.Printf("读取的%d，实际json文件%d字节\n", headSize, len(buf[:headSize]))
		err = json.Unmarshal(buf[:headSize], &head)
		if err != nil {
			conn.Write([]byte{'0'})
			return err
		}

		conn.Write([]byte{'1'}) // 告知对方接收成功
		break
	}

	// 头部接收完毕，生成文件

	//发送文件
	fileName := head.FileList[0]
	err := lib.SendFile(fileName, conn)
	if err != nil {
		return err
	}

	return nil
}
