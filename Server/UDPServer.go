package main

import (
	"TimeSoft-OA/lib"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var udpDirector = map[byte]func(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn){
	byte(lib.Signup):      processSignup,      // 注册
	byte(lib.Login):       processLogin,       // 登录
	byte(lib.ClientCo):    processClientCo,    // 获取客户公司
	byte(lib.WorkLoad):    processWorkload,    // 客户端获取工作量
	byte(lib.EditMission): processEditMission, // 客户端获取工作量
	byte(lib.ReworkItem):  processMarkRework,  // 标记待返工的扫描
}

// UDPServer 用于注册和登录的udp服务器
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

		// log.Printf("收到%s的消息\n", addr.String())
		go processUDPMsg(udpMessage[:n], addr, listen) // 后台处理此条消息

	}
}

// 处理UDP消息
func processUDPMsg(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	packType := udpMessage[0]
	udpMessage = udpMessage[1:]

	fn, ok := udpDirector[packType] // 获取对应函数
	if !ok {
		// 消息有误不做处理
		log.Printf("收到 %s 未识别的消息\n", addr.String())
		return
	}
	fn(udpMessage, addr, listen)
}

// 处理用户注册
func processSignup(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
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
	if err == sql.ErrNoRows { // 未查询到数据
		listen.WriteToUDP(lib.WrongIDOrPwd.Pack(), addr) // 发送错误消息
		log.Printf("%s登录时用户名或密码错误\n", addr.String())
		return
	} else if err != nil {
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		log.Printf("%s登录时操作数据库错误\n", addr.String())
		return
	}

	if loginJson.Pwd != truePwd { // 判断用户名与密码
		listen.WriteToUDP(lib.WrongIDOrPwd.Pack(), addr) // 发送错误消息
		log.Printf("%s登录时用户名或密码错误\n", addr.String())
		return
	}

	// 密码正确
	log.Printf("%s 用户%s登录一次\n", addr.String(), userRealName)
	listen.WriteToUDP(lib.Success.Pack(), addr) // 成功
}

// 处理客户公司名称的获取
func processClientCo(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	var clientCoList []string
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

// 处理客户端查询修图任务
func processEditMission(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {
	var missionList []string

	err := json.Unmarshal(udpMessage, &missionList)
	if err != nil {
		log.Printf("解析消息错误: %v\n", err)
		listen.WriteToUDP(lib.Failed.Pack(), addr) // 发送错误消息
		return
	}

	user := missionList[0]
	missionList = []string{}

	// 获取修图任务数量
	query :=
		`SELECT FILE_ID, CLIENT_COMPANY FROM COMPLETED_TASK
	WHERE EDITOR = ? AND (FILE_STATE = ? OR FILE_STATE = ?);`
	rows, err := db.Query(query, user, lib.ScanOver, lib.ReworkOver)
	if err != nil {
		log.Println("获取扫描任务失败")
		listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		return
	}
	for rows.Next() {
		s1, s2 := "", ""
		err = rows.Scan(&s1, &s2)
		if err != nil {
			log.Println("读取rows失败")
			listen.WriteToUDP(lib.DBOperateErr.Pack(), addr) // 发送错误消息
		}
		missionList = append(missionList, s2+"/"+s1)
	}

	UserEditMap[user] = missionList // 设置全局任务映射，为获取文件做准备
	b, _ := lib.Uint32ToByte(uint32(len(missionList)))

	listen.WriteToUDP(b, addr) // 向客户端返回信息
	log.Printf("地址%s，员工%s查询了一次待修图任务\n", addr, user)
}

// 处理标记返工任务
func processMarkRework(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) {

}
