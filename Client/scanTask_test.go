package main

import (
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
