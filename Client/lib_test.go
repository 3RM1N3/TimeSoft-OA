package main

import (
	"testing"
)

func TestGetMacAddrs(t *testing.T) {
	got, err := GetMacAddrs()
	if err != nil {
		t.Errorf("获取值错误：%v\n", err)
	}

	want := "00:15:5d:90:4d:3d"
	if got[0] != want {
		t.Errorf("\n期待值: %v\n实际值: %v\n", want, got[0])
	}
}
