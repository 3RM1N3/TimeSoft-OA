package main

import (
	"fmt"
	"net"
)

func main() {
	byteChan := make(chan []byte, 128)
	listen, err := net.Listen("tcp", ":8888") // 创建用于监听的 socket
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	fmt.Println("开始监听...")
	go processBytes(byteChan) // 处理收到的字节数据

	defer listen.Close() // 服务器结束前关闭 listener

	// 循环等待客户端来链接
	for {
		fmt.Println("阻塞等待客户端连接...")
		conn, err := listen.Accept() // 创建用户数据通信的socket
		if err != nil {
			fmt.Println("Accept() err=", err)
		} else {
			fmt.Printf("与客户端%s连接成功\n", conn.RemoteAddr().String())
		}

		go process(conn, byteChan) // 起一个协程，为客户端服务
	}

}

func process(conn net.Conn, byteChan chan []byte) {
	firstRead := true
	var byteBuf []byte
	defer conn.Close()
	for {

		buf := make([]byte, 1024) // 创建一个新切片， 用作保存数据的缓冲区
		fmt.Printf("等待%s发送信息...\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf) // 从conn中读取客户端发送的数据内容

		byteBuf = append(byteBuf, buf[:n]...)
		//println(string(byteBuf[len(byteBuf)-6:]))

		if !firstRead && string(byteBuf[len(byteBuf)-6:]) == "\\ioEOF" {
			byteBuf = byteBuf[:len(byteBuf)-6]
			byteChan <- byteBuf // 传出完整[]byte
			//fmt.Println("\n收到了：", string(byteBuf))
			buf = make([]byte, 1024)
			byteBuf = []byte{}
			firstRead = true
		}
		firstRead = false
		if err != nil {
			fmt.Printf("客户端退出，与客户端断开连接\n")
			return
		}
		//fmt.Printf("当前线程 %v, 接受消息 %s\n", goID(), string(buf[:n]))

		// 回写数据给客户端
		_, err = conn.Write([]byte("This is Server"))
		if err != nil {
			fmt.Println("Write err:", err)
			return
		}
	}
}

// 接收处理[]byte
func processBytes(byteChan chan []byte) {
	for {
		println(string(<-byteChan))
	}
}

/*
func goID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
*/
