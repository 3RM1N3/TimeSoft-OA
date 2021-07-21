package SocketPacket

import (
	"bytes"
	"encoding/binary"
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
