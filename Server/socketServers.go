package main

import (
	"TimeSoft-OA/lib"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"
)

// 收发文件的TCP服务器
func TCPServer(address string) error {
	listen, err := net.Listen("tcp", address) // 开启tcp服务器
	if err != nil {
		return err
	}

	log.Printf("TCP服务器在%s开启", address)
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
func UDPServer(address string) {
	//建立一个UDP监听
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Printf("UDP服务未开启, err:%v\n", err)
		os.Exit(1)
		return
	}
	listen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("UDP服务未开启, err:%v\n", err)
		os.Exit(1)
		return
	}

	log.Printf("UDP服务器在%s开启", address)
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
	n, err := conn.Read(b)
	if err != nil {
		log.Println(err)
		return
	} else if n == 0 {
		log.Println("未读取到数据")
		return
	}

	if b[0] == byte(lib.SendHead) { // 客户端发送文件
		head, err := lib.ReceiveFile(conn)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("收到文件", head.Name)

		// 处理文件
		err = StoreFile(head)
		if err != nil {
			log.Println(err)
		}

	} else if b[0] == byte(lib.ReceiveHead) { // 客户端接收文件
		err = ServerSendFile(conn)
		if err != nil {
			log.Println(err)
		}
	}
	// 否则即远端返回了无法识别的消息，函数也返回
}

// 处理UDP消息 待优化
func processUDPMsg(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	packType := udpMessage[0]
	udpMessage = udpMessage[1:]
	switch packType {
	case byte(lib.Signup): // 注册
		processSignin(udpMessage, addr, listen)

	case byte(lib.Login): // 登录
		processLogin(udpMessage, addr, listen)

	case byte(lib.ClientCo): // 获取客户公司
		processClientCo(udpMessage, addr, listen)

	case byte(lib.WorkLoad): // 客户端获取工作量
		processWorkload(udpMessage, addr, listen)

	default:
		// 消息有误不做处理
		log.Printf("收到 %s 未识别的消息\n", addr.String())
	}
}

// 处理用户注册
func processSignin(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	var signup lib.SignUpJson
	err := json.Unmarshal(udpMessage, &signup)
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
	err := json.Unmarshal(udpMessage, &loginJson) // 获取客户端提供的用户名和密码
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

// 处理客户公司名称的获取
func processClientCo(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	clientCoList := []string{}
	rows, err := db.Query(
		`SELECT NAME FROM PROJECT`,
	)
	if err != nil {
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}
	defer rows.Close()

	CoName := ""
	for rows.Next() {
		err := rows.Scan(&CoName)
		if err != nil {
			listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
			return
		}
		clientCoList = append(clientCoList, CoName)
	}

	jsonData, err := json.Marshal(clientCoList)
	if err != nil {
		listen.WriteToUDP(lib.Failed.Pack(), addr) // 发送错误消息
		return
	}

	listen.WriteToUDP(jsonData, addr)
}

// 处理客户端获取工作量
func processWorkload(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	var workloadStruct lib.WorkLoadJson
	var isAdmin = false

	err := json.Unmarshal(udpMessage, &workloadStruct)
	if err != nil {
		log.Printf("解析消息错误: %v\n", err)
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	// 判断是否为管理员
	row := db.QueryRow(`SELECT PHONE FROM ADMIN WHERE PHONE = ?`, workloadStruct.Phone)
	phone := ""
	err = row.Scan(&phone)
	if err == nil {
		isAdmin = true
	} else if err == sql.ErrNoRows {
		isAdmin = false
	} else {
		log.Println("查询管理员表错误")
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	if isAdmin {
		fmt.Println("管理员相关操作")
		return
	}

	// 普通员工
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()

	// 获取扫描任务数量
	query :=
		`SELECT COUNT(FILE_ID) FROM COMPLETED_TASK
	WHERE SCANER = ? AND SCAN_TIME > ? AND FILE_STATE <= 2;`
	row = db.QueryRow(query, workloadStruct.Phone, today)

	err = row.Scan(&workloadStruct.Scan)
	if err != nil {
		log.Println("获取扫描任务失败")
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	// 获取修图任务数量
	query =
		`SELECT COUNT(FILE_ID) FROM COMPLETED_TASK
	WHERE EDITOR = ? AND EDIT_TIME > ? AND FILE_STATE > 1;`
	row = db.QueryRow(query, workloadStruct.Phone, today)

	err = row.Scan(&workloadStruct.Edit)
	if err != nil {
		log.Println("获取修图任务失败")
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	// 获取返工任务数量
	query =
		`SELECT COUNT(FILE_ID) FROM COMPLETED_TASK
	WHERE SCANER = ? AND SCAN_TIME > ? AND FILE_STATE > 2;`
	row = db.QueryRow(query, workloadStruct.Phone, today)

	err = row.Scan(&workloadStruct.Rework)
	if err != nil {
		log.Println("获取返工任务失败")
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}

	jsonData, err := json.Marshal(workloadStruct)
	if err != nil {
		log.Println("生成json失败")
		return
	}

	listen.WriteToUDP(jsonData, addr) // 向客户端返回信息
	log.Printf("地址%s，员工%s查询了一次今日任务量\n", addr, workloadStruct.Phone)
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
	err := lib.SendFile(fileName, "", "", 0x0, conn)
	if err != nil {
		return err
	}

	return nil
}

// 将收到的文件解包归档入库
func StoreFile(head lib.FileSendHead) error {
	fileIDList, err := lib.Unzip(head.Name, path.Join(ProjectPath, head.ClientCo)) // 将压缩文件解压至项目文件夹
	if err != nil {
		return err
	}
	log.Printf("%s发来的文件%s解压完毕\n", head.Uploader, head.Name)

	var stmt *sql.Stmt
	var fileState = 0
	if head.ScanOrEdit == 0x0 { // 扫描
		log.Printf("查询%s等档号是否曾提交过\n", fileIDList[0])

		row := db.QueryRow(`SELECT FILE_STATE FROM COMPLETED_TASK WHERE FILE_ID = ?`, fileIDList[0])
		err = row.Scan(&fileState)
		if err == nil { // 查询到结果
			log.Printf("%s等档号提交过，此次提交为返工\n", fileIDList[0])
			fileState = 5

			stmt, err = db.Prepare(
				`UPDATE COMPLETED_TASK SET CLIENT_COMPANY = ?, SCANER = ?, SCAN_TIME = ?, FILE_STATE = ?, FILE_NUM = ? WHERE FILE_ID = ?`,
			)

		} else if err == sql.ErrNoRows {
			log.Printf("%s等档号未提交过，此次为第一次提交\n", fileIDList[0])
			fileState = 0

			stmt, err = db.Prepare(
				`INSERT INTO COMPLETED_TASK (CLIENT_COMPANY, SCANER, SCAN_TIME, FILE_STATE, FILE_NUM, FILE_ID) VALUES (?, ?, ?, ?, ?, ?)`,
			)

		} else {
			return err
		}

	} else { // 修图
		row := db.QueryRow(`SELECT FILE_STATE FROM COMPLETED_TASK WHERE FILE_ID = ?`, fileIDList[0])
		err = row.Scan(&fileState)
		if err != nil {
			return err
		}
		if fileState == 0 { // 正常扫描待修图
			fileState = 2
		} else if fileState == 5 { // 返工后待修图
			fileState = 6
		} else {
			return errors.New("档案状态有误")
		}
		stmt, err = db.Prepare(
			`INSERT INTO COMPLETED_TASK (CLIENT_COMPANY, EDITOR, EDIT_TIME, FILE_STATE, FILE_NUM, FILE_ID) VALUES (?, ?, ?, ?, ?, ?)`,
		)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, fileID := range fileIDList {
		log.Printf("将%s条目插入数据库\n", fileID)
		subFileNum, err := getSubFileNum(path.Join(head.ClientCo, fileID))
		if err != nil {
			log.Printf("读取项目%s，档号%s内文件数目错误: %v\n", head.ClientCo, fileID, err)
			continue
		}
		_, err = stmt.Exec(head.ClientCo, head.Uploader, time.Now().Unix(), fileState, subFileNum, fileID)
		if err != nil {
			return err
		}
	}

	return nil
}
