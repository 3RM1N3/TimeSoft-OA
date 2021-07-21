package SocketPacket

import (
	"bytes"
	"fmt"
)

type SocketPacket struct {
	PacketType  string
	DataLen     uint32 // uint16
	CurrentPart uint16 // uint16
	AllPart     uint16 // uint16
	Data        []byte
	EndByte     byte
}

func NewSocketPacket(b []byte) SocketPacket {
	return SocketPacket{
		PacketType:  "ZA",
		DataLen:     uint32(len(b) + 9),
		CurrentPart: uint16(1),
		AllPart:     uint16(1),
		Data:        b,
		EndByte:     0x01,
	}
}

func (sp *SocketPacket) String() string {
	return fmt.Sprintf("Type: SocketPacket\n    PacketType:  %s,\n    DataLen:     %d,\n    CurrentPart: %d,\n    AllPart:     %d,\n    Data:        []byte{...},\n    EndByte:     %X,",
		sp.PacketType,
		sp.DataLen,
		sp.CurrentPart,
		sp.AllPart,
		sp.EndByte)
}

func (sp *SocketPacket) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte(sp.PacketType)) // 写入类型

	get2bytes, _ := Uint32ToByte(sp.DataLen)
	buf.Write(get2bytes)

	get2bytes, _ = Uint16ToByte(sp.CurrentPart)
	buf.Write(get2bytes)

	get2bytes, _ = Uint16ToByte(sp.AllPart)
	buf.Write(get2bytes)

	buf.Write(sp.Data)
	buf.WriteByte(sp.EndByte)

	return buf.Bytes()
}
