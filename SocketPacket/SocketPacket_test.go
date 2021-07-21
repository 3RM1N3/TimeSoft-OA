package SocketPacket

import (
	"fmt"
	"testing"
)

func TestNewSocketPacket(t *testing.T) {
	got := NewSocketPacket([]byte{0x00, 0x00, 0xf0, 0xfa})

	want := uint16(61690)
	if got.PacketType != "AA" {
		t.Errorf("expect %v, however %v", want, got)
	}
}

func ExampleNewSocketPacket() {
	got := NewSocketPacket([]byte{0x00, 0x00, 0xf0, 0xfa})
	fmt.Println(got)
}
