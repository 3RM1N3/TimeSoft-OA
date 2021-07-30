package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"testing"
)

func TestByteToUint16(t *testing.T) {
	got, err := SP.ByteToUint16([]byte{0x00, 0x00, 0xf0, 0xfa})
	if err != nil {
		t.Errorf("转换失败: %v", err)
		return
	}
	want := uint16(61690)
	if got != want {
		t.Errorf("expect %v, however %v", want, got)
	}
}
