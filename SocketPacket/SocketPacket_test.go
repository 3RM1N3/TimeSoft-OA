package SocketPacket

import (
	"testing"
)

func TestNewSocketPacket(t *testing.T) {
	ch := make(chan SocketPacket)
	got := ""
	NewJsonPacket(Login, []byte{0x00, 0x00, 0xf0, 0xfa}, ch)

	want := uint16(61690)
	if got != "AA" {
		t.Errorf("expect %v, however %v", want, got)
	}
}
