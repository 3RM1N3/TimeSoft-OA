package SocketPacket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

type SocketPacket struct {
	TypeByte    PacketType
	DataLen     uint32
	CurrentPart uint16
	AllPart     uint16
	Data        []byte
}

// 显示此SocketPacket
func (sp *SocketPacket) String() string {
	return fmt.Sprintf("Type: SocketPacket\n    PacketType:  %X,\n    DataLen:     %d,\n    CurrentPart: %d,\n    AllPart:     %d,\n    Data:        []byte{...}",
		sp.TypeByte,
		sp.DataLen,
		sp.CurrentPart,
		sp.AllPart,
	)
}

// 将SocketPacket内容打包为字节切片。格式：
//
// var pack []byte
//
// pack[0] 文件类型
//
// pack[1:5] 此包数据正文长度        最大值 8 << 20
//
// pack[5:7] 此包为该文件的第几个包   最大值512
//
// pack[7:9] 此文件一共有多少包      最大值512
//
// pack[9:len(pack)-1] 数据正文
//
// pack[len(pack)-1] 结束标识符     0x00 / 0x01
func (sp *SocketPacket) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(byte(sp.TypeByte)) // 写入类型
	if err != nil {
		return []byte{}, err
	}

	get2bytes, err := Uint32ToByte(sp.DataLen)
	if err != nil {
		return []byte{}, err
	}
	buf.Write(get2bytes)

	get2bytes, err = Uint16ToByte(sp.CurrentPart)
	if err != nil {
		return nil, err
	}
	buf.Write(get2bytes)

	get2bytes, err = Uint16ToByte(sp.AllPart)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(get2bytes)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(sp.Data)
	if err != nil {
		return nil, err
	}

	difference := packageSize - buf.Len() // 补齐packageSize个字节

	if difference < 0 { // 如果长度大于一个包，返回错误
		return nil, errors.New("传入数据过长")

	} else if difference > 0 { // 如果长度不够一个包则补齐
		buf.Write(make([]byte, difference))
	}

	return buf.Bytes(), nil
}

// 从字节缓冲区读取包数据
func (sp *SocketPacket) ReadPack(b *[]byte) error {
	var err error

	sp.TypeByte = PacketType((*b)[0]) // 获取类型

	sp.DataLen, err = ByteToUint32((*b)[1:5]) // 获取数据长度
	if err != nil {
		return err
	}

	sp.CurrentPart, err = ByteToUint16((*b)[5:7]) // 获取当前包索引
	if err != nil {
		return err
	}

	sp.AllPart, err = ByteToUint16((*b)[7:9]) // 获取全部包数量
	if err != nil {
		return err
	}

	sp.Data = (*b)[9 : 9+int(sp.DataLen)] // 获取正文数据
	(*b) = (*b)[packageSize:]             // 丢弃空白部分

	return nil
}

// 读取字节切片生成json packet包传入out
func NewJsonPacket(t PacketType, b []byte, out chan SocketPacket) {
	out <- SocketPacket{
		TypeByte:    t,
		DataLen:     uint32(len(b)),
		CurrentPart: 1,
		AllPart:     1,
		Data:        b,
	}
}

// 打包zip数据并拆分成包传入out
func NewZipPacket(f *os.File, out chan SocketPacket) error {
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	zipSize := stat.Size()
	allPart := math.Ceil(float64(zipSize) / packageSize)

	fileJson := FileUploadJson{
		FileName: stat.Name(),
		FileSize: zipSize,
	}

	encodedFileJson, err := json.Marshal(fileJson)
	if err != nil {
		return err
	}

	out <- SocketPacket{ // 传入上传文件json数据
		TypeByte:    FileUpload,
		DataLen:     uint32(len(encodedFileJson)),
		CurrentPart: 0,
		AllPart:     uint16(allPart),
		Data:        encodedFileJson,
	}

	for i := 0; i < int(allPart); i++ {
		buf := make([]byte, packageSize-9)
		n, err := f.Read(buf)
		if err != nil {
			return err
		}

		out <- SocketPacket{ // 循环传送文件包
			TypeByte:    ZipArchive,
			DataLen:     uint32(n),
			CurrentPart: uint16(i) + 1,
			AllPart:     uint16(allPart),
			Data:        buf[:n],
		}
	}
	return nil
}
