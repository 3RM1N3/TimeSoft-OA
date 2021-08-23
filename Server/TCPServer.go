package main

import (
	"TimeSoft-OA/lib"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

// TCPServer 收发文件的TCP服务器
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
		log.Printf("客户端%s请求下载文件\n", conn.RemoteAddr().String())
		err = ServerSendFile(conn)
		if err != nil {
			log.Println(err)
		}
	}
	// 否则即远端发送的为无法识别的消息，函数返回
}

// ServerSendFile TCP服务器发送文件给远端
func ServerSendFile(conn net.Conn) error {
	var head lib.FileReceiveHead
	headSize := uint16(0)
	var buf []byte

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
	fileName := head.Downloader + ".temp"
	err := ZipDirs(UserEditMap[head.Downloader], fileName)
	if err != nil {
		return fmt.Errorf("生成压缩文件失败：%v", err)
	}

	//发送文件
	err = lib.SendFile(fileName, "", "", 0x0, conn)
	if err != nil {
		return err
	}

	return nil
}


// StoreFile 将收到的文件解包归档入库
func StoreFile(head lib.FileSendHead) error {
	defer os.Remove(head.Name) // 执行完毕后删除文件

	// 解压文件
	fileIDList, err := lib.Unzip(head.Name, filepath.Join(ProjectPath, head.ClientCo)) // 将压缩文件解压至项目文件夹
	if err != nil {
		return err
	}
	log.Printf("%s发来的文件%s解压完毕\n", head.Uploader, head.Name)

	// 分情况操作数据库
	if head.ScanOrEdit == 0x0 { // 扫描
		if head.IsRework {
			err = ReceiveScanRework(head, fileIDList)
		} else {
			err = ReceiveScanFile(head, fileIDList)
		}
		if err != nil {
			return err
		}
	} else { // 修图
		if head.IsRework {
			err = ReceiveEditRework(head, fileIDList)
		} else {
			err = ReceiveEditFile(head, fileIDList)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// ReceiveScanFile 接收扫描的文件提交
func ReceiveScanFile(head lib.FileSendHead, fileIDList []string) error {
	var query string
	insertQuery := `INSERT INTO COMPLETED_TASK
                  (CLIENT_COMPANY, SCANER, SCAN_TIME, FILE_STATE, FILE_NUM, FILE_ID)
    			  VALUES (?, ?, ?, ?, ?, ?);`
	updateQuery := `UPDATE COMPLETED_TASK
                  SET CLIENT_COMPANY = ?, SCANER = ?, SCAN_TIME = ?, FILE_STATE = ?, FILE_NUM = ?
    			  WHERE FILE_ID = ?;`
	existQuery := `SELECT COUNT(FILE_ID) FROM COMPLETED_TASK WHERE FILE_ID = ?`

	for _, j := range fileIDList {
		// 查询是否为第一次提交
		row := db.QueryRow(existQuery, j)
		var c int
		if err := row.Scan(&c); err != nil {
			return fmt.Errorf("查询数据库错误：%v", err)
		}
		if c == 0 {
			query = insertQuery
		} else {
			query = updateQuery
		}

		// 写入数据库信息
		n, err := getSubFileNum(filepath.Join(ProjectPath, head.ClientCo, j))
		if err != nil {
			return fmt.Errorf("读取文件数错误：%v", err)
		}
		_, err = db.Exec(query, head.ClientCo, head.Uploader, time.Now().Unix(), lib.ScanOver, n, j)
		if err != nil {
			return fmt.Errorf("操作数据库错误：%v", err)
		}
	}

	return nil
}

// ReceiveScanRework 接收扫描返工的文件提交
func ReceiveScanRework(head lib.FileSendHead, fileIDList []string) error {
	var query = `UPDATE COMPLETED_TASK
                  SET CLIENT_COMPANY = ?, SCANER = ?, SCAN_TIME = ?, FILE_STATE = ?, FILE_NUM = ?
    			  WHERE FILE_ID = ?;`
	for _, j := range fileIDList {

		// 写入数据库信息
		n, err := getSubFileNum(filepath.Join(ProjectPath, head.ClientCo, j))
		if err != nil {
			return fmt.Errorf("读取文件数错误：%v", err)
		}
		_, err = db.Exec(query, head.ClientCo, head.Uploader, time.Now().Unix(), lib.ScanOver, n, j)
		if err != nil {
			return fmt.Errorf("操作数据库错误：%v", err)
		}
	}

	return nil
}

// ReceiveEditFile 接收修图的文件提交
func ReceiveEditFile(head lib.FileSendHead, fileIDList []string) error {
	// 写入数据库信息
	n, err := getSubFileNum(filepath.Join(ProjectPath, head.ClientCo, j))
	if err != nil {
		return fmt.Errorf("读取文件数错误：%v", err)
	}
	_, err = db.Exec(query, head.ClientCo, head.Uploader, time.Now().Unix(), lib.EditOver, n, j)
	if err != nil {
		return fmt.Errorf("操作数据库错误：%v", err)
	}
}

// ReceiveEditRework 接收修图返工的文件提交
func ReceiveEditRework(head lib.FileSendHead, fileIDList []string) error {
	// 写入数据库信息
	n, err := getSubFileNum(filepath.Join(ProjectPath, head.ClientCo, j))
	if err != nil {
		return fmt.Errorf("读取文件数错误：%v", err)
	}
	_, err = db.Exec(query, head.ClientCo, head.Uploader, time.Now().Unix(), lib.EditOver, n, j)
	if err != nil {
		return fmt.Errorf("操作数据库错误：%v", err)
	}
}
