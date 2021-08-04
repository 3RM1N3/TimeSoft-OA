package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var mapIPuser = map[string]string{}

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

		if bytePack[0] == byte(SP.Signup) { // 注册
			var signup SP.SignUpJson
			err := json.Unmarshal(bytePack[1:n], &signup)
			if err != nil {
				reportSP.Msg = "发送的数据有误"
			} else {
				log.Printf("%s 申请注册\n", addr.String())
				// 判断与数据库中的项目有无重复
				// 从数据库查询对应密码

				row := db.QueryRow(`SELECT PHONE FROM USER WHERE PHONE = ? UNION SELECT PHONE FROM UNREVIEWED_USER WHERE PHONE = ?`, signup.PhoneNumber, signup.PhoneNumber)
				err := row.Scan()
				if err != sql.ErrNoRows { // 如果有数据
					reportSP.Msg = "此电话号码已被注册"
				} else { // 在已有用户中唯一，可以加入未审核表中
					_, err := db.Exec(`INSERT INTO UNREVIEWED_USER VALUES (?, ?, ?)`, signup.PhoneNumber, signup.Pwd, signup.RealName)
					if err != nil {
						reportSP.Msg = "写入数据库错误"
					} else {
						reportSP.Success = true
						fmt.Printf("新员工注册待审核：\n    电话号码：%s\n    真实姓名：%s\n", signup.PhoneNumber, signup.RealName)
					}
				}
			}

		} else if bytePack[0] == byte(SP.Login) { // 登录
			userRealName := ""
			loginJson := SP.LoginJson{}
			err := json.Unmarshal(bytePack[1:n], &loginJson) // 获取客户端提供的用户名和密码
			if err != nil {
				reportSP.Msg = "发送的数据有误"
			} else {
				// 从数据库查询对应密码
				row := db.QueryRow(`SELECT PWD, REALNAME FROM USER WHERE PHONE = ?`, loginJson.PhoneNumber)
				truePwd := ""
				err = row.Scan(&truePwd, &userRealName)
				if err != nil {
					reportSP.Msg = "用户名或密码错误"
				} else if loginJson.Pwd != truePwd { // 判断用户名与密码
					reportSP.Msg = "用户名或密码错误"
				} else { // 密码正确
					mapIPuser[addr.String()] = userRealName + loginJson.PhoneNumber[7:]
					log.Printf("%s 用户%s登录成功\n", addr.String(), userRealName)
					reportSP.Success = true
				}
			}

		} else { // 消息有误
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
	userRealName := ""
	var fileSize, writeSize int64

	for {
		sp := <-receiveSPChan

		switch sp.TypeByte {

		// 客户端上传文件信息至服务器
		case SP.FileUpload:
			if userRealName == "" {
				log.Println("用户未登录，但尝试上传文件")
				continue
			}
			fileUploadJson := SP.FileUploadJson{}
			json.Unmarshal(sp.Data, &fileUploadJson)
			fileName = fileUploadJson.FileName
			fileSize = fileUploadJson.FileSize
			fmt.Printf("\n收到文件%s，大小%d\n", fileName, fileSize)
			writeSize = 0

		// 保存文件
		case SP.ZipArchive:
			if userRealName == "" {
				continue
			}

			// 在本地创建文件
			fmt.Printf("\n***此次写入文件 %s 的第%d部分***\n", fileName, sp.CurrentPart)
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				log.Println("文件打开失败", err)
				//continue
				return
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

// 获取未审核员工信息
func GetUnreviewUser(unreviewedUserList *[][2]string) error {
	rows, err := db.Query(`SELECT PHONE, REALNAME FROM UNREVIEWED_USER`)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Println("查询数据库失败")
		return err
	}

	for rows.Next() {
		tempList := [2]string{}
		err = rows.Scan(&tempList[0], &tempList[1])
		if err != nil {
			log.Println("获取某条内容失败")
			return err
		}
		*unreviewedUserList = append(*unreviewedUserList, tempList)
	}

	return nil
}

// 令未审核的员工全部通过
func AllPass() error {
	_, err := db.Exec(`INSERT INTO USER SELECT * FROM UNREVIEWED_USER`)
	if err != nil {
		log.Println("未审核员工通过失败")
		return err
	}
	_, err = db.Exec(`DELETE FROM UNREVIEWED_USER`)
	if err != nil {
		log.Println("清空未审核表失败，请手动清空")
		return err
	}

	return nil
}

// 部分通过
func PartPass(unreviewedList *[][2]string, passedIndex *[]int) error {
	for i, v := range *unreviewedList {
		if intListContains(*passedIndex, i) {
			// 通过用户插入至用户表
			_, err := db.Exec(`INSERT INTO USER SELECT * FROM UNREVIEWED_USER WHERE PHONE = ?`, v[0])
			if err != nil {
				log.Println("插入失败")
				return err
			}
		}
		// 删除已审核记录
		_, err := db.Exec(`DELETE FROM UNREVIEWED_USER WHERE PHONE = ?`, v[0])
		if err != nil {
			log.Println("UNREVIEWED_USER中删除此条记录失败，请手动删除 电话号码:", v[0])
			return err
		}
	}
	return nil
}

func intListContains(l []int, i int) bool {
	for _, e := range l {
		if e == i {
			return true
		}
	}
	return false
}

// 重置密码
func ResetPwd(phone, newPwd string) error {
	md5pwd := SP.MD5(newPwd)

	_, err := db.Exec(`UPDATE USER SET PWD = ? WHERE PHONE = ?`, md5pwd, phone)
	if err != nil {
		log.Println("重置失败")
		return err
	}

	return nil
}
