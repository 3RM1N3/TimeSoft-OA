package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

// 上传文件的请求头
type FileSendHead struct {
	Name     string
	Uploader string
	Size     int64
}

// 下载文件的请求头
type FileReceiveHead struct {
	FileList   []string
	Downloader string
}

func (h *FileSendHead) MakeHead() ([]byte, error) {
	return MakeHead(SendHead, h)
}

func (h FileReceiveHead) MakeHead() ([]byte, error) {
	return MakeHead(ReceiveHead, h)
}

// 创建请求头字节切片
func MakeHead(SRType PacketType, some interface{}) ([]byte, error) {
	b, err := json.Marshal(some)
	if err != nil {
		return nil, err
	}

	headByte, err := Uint16ToByte(uint16(len(b)))
	if err != nil {
		return nil, err
	}
	headByte = append([]byte{byte(SRType)}, headByte...)

	return append(headByte, b...), nil
}

// 发送文件至远端
func SendFile(fileName string, conn net.Conn) error {
	fmt.Println("准备发送文件", fileName)
	defer conn.Close()

	state, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	partNum := math.Ceil(float64(state.Size()) / 8388608)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	uploadHead := FileSendHead{
		Name:     fileName,
		Size:     state.Size(),
		Uploader: "13284030601",
	}
	fileHead, err := uploadHead.MakeHead()
	if err != nil {
		return err
	}

	_, err = conn.Write(fileHead)
	if err != nil {
		return err
	}
	fmt.Println("发送文件头成功")

	result := make([]byte, 1)
	conn.Read(result)
	if result[0] != '1' {
		return errors.New("远端接收文件时出现错误")
	}
	fmt.Println("可以发送正文")

	filePart := make([]byte, 8388608)
	for i := float64(0); i < partNum; {
		fmt.Printf("第%d次发送\n", int(i)+1)

		n, err := f.Read(filePart)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = conn.Write(filePart[:n])
		if err != nil {
			return err
		}

		_, err = conn.Read(result)
		if err != nil {
			return err
		}

		if result[0] != '1' {
			return errors.New("服务器拒绝接收文件")
		}

		i++
	}

	fmt.Println("文件发送成功")
	return nil
}

// 从远端接收文件，返回本地文件名，文件指针和错误
func ReceiveFile(conn net.Conn) (string, *os.File, error) {
	var f *os.File
	var head FileSendHead
	writtenSize := 0
	headSize := uint16(0)
	buf := []byte{}

	for {
		b := make([]byte, 524288)
		n, err := conn.Read(b)
		if err != nil {
			return "", nil, err
		}

		buf = append(buf, b[:n]...)
		if len(buf) < 3 {
			continue
		}

		if headSize == 0 {
			headSize, err = ByteToUint16(buf[:2])
			if err != nil {
				return "", nil, err
			}
			buf = buf[2:]
		}

		if f == nil {
			if len(buf) < int(headSize) {
				continue
			}

			// 解析头部
			fmt.Printf("读取的%d，实际json文件%d字节\n", headSize, len(buf[:headSize]))
			err := json.Unmarshal(buf[:headSize], &head)
			if err != nil {
				conn.Write([]byte{'0'})
				return "", nil, err
			}
			// 创建本地文件
			f, err = os.OpenFile(head.Name, os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				conn.Write([]byte{'0'})
				return "", nil, err
			}
			// defer f.Close()
			fmt.Println("创建本地文件成功")

			conn.Write([]byte{'1'}) // 告知对方接收成功

			buf = []byte{}
			continue

			//buf = buf[headSize:] // 数据多于头部的情况不可能发生
		}

		// 存储文件
		n, err = f.Write(buf)
		if err != nil {
			return "", nil, err
		}
		writtenSize += n
		buf = []byte{}
		fmt.Println("此次写入", n, "字节")

		if writtenSize%8388608 == 0 {
			conn.Write([]byte{'1'})
			fmt.Println("收到一个包")
		}

		if writtenSize == int(head.Size) {
			break
		}
	}
	// 接收完毕
	//处理文件
	fmt.Printf("文件%s接收成功\n", head.Name)
	_, err := f.Seek(0, 0)
	if err != nil {
		f.Close()
		return head.Name, nil, err
	}
	return head.Name, f, err
}
