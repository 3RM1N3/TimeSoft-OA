package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	SP "TimeSoft-OA/SocketPacket"
)

func main() {
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
		fmt.Printf("与客户端%s连接成功\n", conn.RemoteAddr().String())
		go process(conn) // 起一个协程，为该客户端服务
	}

}

func process(conn net.Conn) {
	receivedSPChan := make(chan SP.SocketPacket, 128)
	sendSPChan := make(chan SP.SocketPacket, 64)

	go send(conn, sendSPChan)       // 发送数据
	go processBytes(receivedSPChan) // 处理收到的字节数据

	read(conn, receivedSPChan, sendSPChan)
}

// 读取部分
func read(conn net.Conn, receivedSPChan chan SP.SocketPacket, sendSPChan chan SP.SocketPacket) {
	firstRead := true
	byteBuf := new(bytes.Buffer)
	var readedLength uint32
	var sp SP.SocketPacket
	defer conn.Close()

	for {
		buf := make([]byte, 1024) // 创建一个新切片， 用作保存数据的缓冲区
		fmt.Printf("\n等待%s发送信息...\n", conn.RemoteAddr().String())

		n, err := conn.Read(buf) // 读取数据，无则阻塞
		if err != nil {
			log.Printf("客户端退出，与%s断开连接\n", conn.RemoteAddr().String())
			return
		}

		println(conn.RemoteAddr().String(), "此批接收完毕，读取字符数：", n)
		println("内容：", string(buf[:n]))
		readedLength += uint32(n)
		byteBuf.Write(buf[:n])

		if firstRead {
			println(conn.RemoteAddr().String(), "第一次读取，判断文件头")
			sp, _ = ProcessPackHead(byteBuf)
			println(conn.RemoteAddr().String(), "此次接收类型：", sp.TypeByte, "  长度：", sp.DataLen)
			firstRead = false
		}

		if sp.DataLen <= readedLength {
			println(conn.RemoteAddr().String(), "此包传输完毕待处理")
			i := sp.DataLen + uint32(n) - readedLength - 1
			fmt.Printf("%s 终止符位置在：%d，计算得出的终止符为：%x\n", conn.RemoteAddr().String(), i, buf[i])
			if buf[i] == 0x01 {
				println(conn.RemoteAddr().String(), "终止符相同，处理和传出数据")
				sp.Data = byteBuf.Next(int(sp.DataLen) - 10)
				fmt.Printf("传出数据长度为：%d\n", len(sp.Data))
				fmt.Println("传出数据内容为：\n", string(sp.Data))
				readedLength = 0
				byteBuf.ReadByte()
				firstRead = true
				//byteChan <- sp
			}
		}
	}

	// 返回数据
	report := ReportJson{
		Success: true,
		Msg:     "",
	}
	reportJsonByte, _ := json.Marshal(report)
	SP.NewJosnPacket(SP.Report, reportJsonByte, sendSPChan)
}

// 处理数据头
func ProcessPackHead(byteBuf *bytes.Buffer) (SP.SocketPacket, error) {
	sp := SP.SocketPacket{}
	if byteBuf.Len() < 9 {
		return sp, errors.New("[]byte is too short")
	}
	typeByte, _ := byteBuf.ReadByte()
	sp.TypeByte = SP.PacketType(typeByte)
	sp.DataLen, _ = SP.ByteToUint32(byteBuf.Next(4))
	sp.CurrentPart, _ = SP.ByteToUint16(byteBuf.Next(2))
	sp.AllPart, _ = SP.ByteToUint16(byteBuf.Next(2))

	return sp, nil
}

// 向客户端发送数据
func send(conn net.Conn, spChan chan SP.SocketPacket) error {
	sp := <-spChan
	b := sp.Pack()
	buf := make([]byte, 32)

	breader := bytes.NewReader(b)

	for {
		n, err := breader.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if _, err := conn.Write(buf[:n]); err != nil {
			return err
		}
	}
}

// 接收处理[]byte
func processBytes(byteChan chan SP.SocketPacket) {
	sp := <-byteChan
	println(string(sp.Data))
}
