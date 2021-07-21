package SocketPacket

import (
	"bytes"
	"encoding/binary"
)

type PacketType byte
type PacketEndByte byte

type FileUploadJson struct {
	FileName string
	FileSize int64
}

type ReportJson struct {
	Success bool
	Msg     string
}

const (
	ZipArchive  PacketType = iota // 压缩文件档案
	Report                        // 汇报，用于汇报数据接受情况
	FileUpload                    // 上传文件之前的汇报
	FileRequest                   // 请求下载文件
	Login                         // 用于登录
	Notice                        // 用于通知
	PullList                      // 客户端请求获得数据表
	PushList                      // 服务端用于推送表
)

const (
	NotOverEnd PacketEndByte = iota // 此包结束但文件未结束
	OverEnd                         // 此包结束且文件已结束
)

func Uint16ToByte(i uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToUint16(b []byte) (uint16, error) {
	if len(b) == 0 {
		return 0, nil
	}
	i := uint16(0)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func Uint32ToByte(i uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToUint32(b []byte) (uint32, error) {
	if len(b) == 0 {
		return 0, nil
	}
	i := uint32(0)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}
