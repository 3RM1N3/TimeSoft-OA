package main

import (
	"errors"
	"testing"
)

func TestCheckMD5(t *testing.T) {
	got := CheckMD5("admin")
	want := "21232f297a57a5a743894a0e4a801fc3"
	if got != want {
		t.Errorf("expect %v, however %v", want, got)
	}
}

func TestZip(t *testing.T) {
	got := Zip(`D:\projects\computer-exam`, "a.zip")
	want := errors.New("")
	if got != want {
		t.Errorf("expect %v, however %v", want, got)
	}
}
