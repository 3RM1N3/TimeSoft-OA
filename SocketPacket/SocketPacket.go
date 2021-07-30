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
	EndByte     PacketEndByte
}

// 显示此SocketPacket
func (sp *SocketPacket) String() string {
	return fmt.Sprintf("Type: SocketPacket\n    PacketType:  %X,\n    DataLen:     %d,\n    CurrentPart: %d,\n    AllPart:     %d,\n    Data:        []byte{...},\n    EndByte:     %X,",
		sp.TypeByte,
		sp.DataLen,
		sp.CurrentPart,
		sp.AllPart,
		sp.EndByte)
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
	err = buf.WriteByte(byte(sp.EndByte))
	if err != nil {
		return nil, err
	}

	// 补齐8388608个字节
	standardLen := 8388608
	difference := standardLen - buf.Len()

	if difference < 0 {
		return nil, err

	} else if difference > 0 {
		buf.Write(make([]byte, difference))
	}

	return buf.Bytes(), nil
}

// 从字节缓冲区读取包数据
func (sp *SocketPacket) ReadPack(b []byte) error {
	if len(b) < 9 {
		return errors.New("[]byte is too short")
	}

	sp.TypeByte = PacketType(b[0])

	var err error
	sp.DataLen, err = ByteToUint32(b[1:5])
	if err != nil {
		return err
	}

	sp.CurrentPart, err = ByteToUint16(b[5:7])
	if err != nil {
		return err
	}

	sp.AllPart, err = ByteToUint16(b[7:9])
	if err != nil {
		return err
	}

	sp.Data = b[9 : 9+sp.DataLen]
	sp.EndByte = PacketEndByte(b[9+sp.DataLen])

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
		EndByte:     OverEnd,
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
	allPart := math.Ceil(float64(zipSize) / 8388598)

	fileJson := FileUploadJson{
		FileName: stat.Name(),
		FileSize: zipSize,
	}

	encodedFileJson, err := json.Marshal(fileJson)
	if err != nil {
		return err
	}

	out <- SocketPacket{
		TypeByte:    FileUpload,
		DataLen:     uint32(len(encodedFileJson)),
		CurrentPart: 0,
		AllPart:     uint16(allPart),
		Data:        encodedFileJson,
		EndByte:     NotOverEnd,
	}

	endByte := NotOverEnd

	for i := 0; i < int(allPart); i++ {
		buf := make([]byte, 8388598)
		n, err := f.Read(buf)
		if err != nil {
			return err
		}

		if i == int(allPart)-1 {
			endByte = OverEnd
		}

		out <- SocketPacket{
			TypeByte:    ZipArchive,
			DataLen:     uint32(n),
			CurrentPart: uint16(i) + 1,
			AllPart:     uint16(allPart),
			Data:        buf[:n],
			EndByte:     endByte,
		}
	}
	return nil
}

// 从in读取包，合并为分开前的文件写入磁盘，返回*os.File
//
// 注意：应仅当SocketPacket.TypeByte == FileUpload时使用此方法保存文件
func (sp *SocketPacket) SplicingFile(spIn chan SocketPacket) error {
	if sp.TypeByte != FileUpload {
		return errors.New("not a file upload json")
	}

	thisFile := FileUploadJson{}
	if err := json.Unmarshal(sp.Data, &thisFile); err != nil {
		return err
	}

	f, err := os.Create(thisFile.FileName)
	if err != nil {
		return err
	}
	defer f.Close()

	sumByte := 0
	for i := 0; i < int(sp.AllPart); i++ {
		spFile := <-spIn
		if spFile.TypeByte != ZipArchive || spFile.CurrentPart != uint16(i)+1 {
			return errors.New("file receive abort")
		}
		//fmt.Printf("\n此包内容：\n%v\n", spFile.Data)
		n, err := f.Write(spFile.Data)
		if err != nil {
			return err
		}
		sumByte += n
		fmt.Printf("写入了%d字节\n", sumByte)
		err = f.Sync()
		if err != nil {
			return err
		}
		fmt.Printf("文件第%d部分写入成功\n", i+1)
	}
	return nil
}
