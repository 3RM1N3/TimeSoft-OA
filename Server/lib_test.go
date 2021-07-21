package main

import (
	"bytes"
	"testing"
)

func TestByteToUint16(t *testing.T) {
	got, err := ByteToUint16([]byte{0x00, 0x00, 0xf0, 0xfa})
	if err != nil {
		t.Errorf("转换失败: %v", err)
		return
	}
	want := uint16(61690)
	if got != want {
		t.Errorf("expect %v, however %v", want, got)
	}
}

func TestUint16ToByte(t *testing.T) {
	got, err := Uint16ToByte(uint16(61690))
	if err != nil {
		t.Errorf("转换失败")
		return
	}
	want := []byte{0xf0, 0xfa}
	if !bytes.Equal(got, want) {
		t.Errorf("expect %v, however %v", want, got)
	}
}
