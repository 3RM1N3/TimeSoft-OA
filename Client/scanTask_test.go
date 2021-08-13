package main

import (
	"TimeSoft-OA/lib"
	"testing"
)

func TestGenVerf(t *testing.T) {
	got, err := GenVerf("./testDir", -1)
	if err != nil {
		t.Errorf("生成校验码错误：%v\n", err)
	}

	want := "verify!#3118717729B11AF36609AD4E99935190"
	if got != want {
		t.Errorf("\n期待值: %v\n实际值: %v\n", want, got)
	}
}

func TestCheckVerf(t *testing.T) {
	ok, err := CheckVerf("./testDir")
	if err != nil {
		t.Errorf("校验发生错误：%v\n", err)
	}

	if !ok {
		t.Errorf("\n文件夹校验未通过\n")
	}
}

func TestSetProjectDir(t *testing.T) {
	err := SetProjectDir("./testDir")
	if err != nil {
		t.Errorf("设置失败：%v\n", err)
	}
}

func TestSetProjectDir1(t *testing.T) {
	err := SetProjectDir("./testProject")
	if err != nil {
		t.Errorf("设置失败：%v\n", err)
	}
}

func TestSetProjectDir2(t *testing.T) {
	err := SetProjectDir("./rfevgrftebgtbv")
	if err == nil {
		t.Errorf("设置成功，但文件夹./rfevgrftebgtbv不应通过测试\n")
	}
}

func TestScanOverPackSubmit(t *testing.T) {
	globalServerAddr = "127.0.0.1"
	a := false
	err := ScanOverPackSubmit("./testDir", "中石油", "13284030601", 0x0, &a)
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestTodayWorkload(t *testing.T) {
	globalServerAddr = "127.0.0.1"
	globalPhone = "13284030601"
	got, err := TodayWorkload()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	want := lib.WorkLoadJson{
		Phone:  globalPhone,
		Scan:   3,
		Edit:   0,
		Rework: 1,
	}

	if got.Scan != want.Scan || got.Edit != want.Edit || got.Rework != want.Rework {
		t.Errorf("期待：%v\n实际：%v\n", want, got)
	}
}
