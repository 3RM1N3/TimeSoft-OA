package SocketPacket

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func NewPacket(b []byte) SocketPacket {
	return SocketPacket{
		TypeByte:    ZipArchive,
		DataLen:     uint32(len(b) + 10),
		CurrentPart: uint16(1),
		AllPart:     uint16(1),
		Data:        b,
		EndByte:     OverEnd,
	}
}

func (sp *SocketPacket) String() string {
	return fmt.Sprintf("Type: SocketPacket\n    PacketType:  %X,\n    DataLen:     %d,\n    CurrentPart: %d,\n    AllPart:     %d,\n    Data:        []byte{...},\n    EndByte:     %X,",
		sp.TypeByte,
		sp.DataLen,
		sp.CurrentPart,
		sp.AllPart,
		sp.EndByte)
}

func (sp *SocketPacket) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(sp.TypeByte)) // 写入类型

	get2bytes, _ := Uint32ToByte(sp.DataLen)
	buf.Write(get2bytes)

	get2bytes, _ = Uint16ToByte(sp.CurrentPart)
	buf.Write(get2bytes)

	get2bytes, _ = Uint16ToByte(sp.AllPart)
	buf.Write(get2bytes)

	buf.Write(sp.Data)
	buf.WriteByte(byte(sp.EndByte))

	return buf.Bytes()
}

func NewJosnPacket(t PacketType, b []byte, out chan SocketPacket) {
	out <- SocketPacket{
		TypeByte:    t,
		DataLen:     uint32(len(b) + 10),
		CurrentPart: 1,
		AllPart:     1,
		Data:        b,
		EndByte:     OverEnd,
	}
}

// 打包zip数据并拆分成包传入out
func NewZipPacket(f *os.File, out chan SocketPacket) error {
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	zipSize := stat.Size()
	allPart := zipSize/8388608 + 1

	fileJson := FileUploadJson{
		FileName: stat.Name(),
		FileSize: zipSize,
	}

	encodedFileJson, _ := json.Marshal(fileJson)

	out <- SocketPacket{
		TypeByte:    FileUpload,
		DataLen:     uint32(len(encodedFileJson) + 10),
		CurrentPart: 0,
		AllPart:     uint16(allPart),
		Data:        encodedFileJson,
		EndByte:     NotOverEnd,
	}

	buf := make([]byte, 8388608)
	endByte := NotOverEnd

	for i := 0; i < int(allPart); i++ {
		n, err := f.Read(buf)
		if err != nil {
			return err
		}

		if i == int(allPart)-1 {
			endByte = OverEnd
		}

		out <- SocketPacket{
			TypeByte:    ZipArchive,
			DataLen:     uint32(n + 10),
			CurrentPart: uint16(i) + 1,
			AllPart:     uint16(allPart),
			Data:        buf[:n],
			EndByte:     endByte,
		}
	}
	return nil
}
