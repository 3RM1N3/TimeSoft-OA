package main

import (
	"testing"
)

// 测试获取未审核用户
func TestGetUnreviewUser(t *testing.T) {
	want := [][2]string{
		{"15146209355", "李晓"},
		{"18512341234", "李雷"},
	}

	got := [][2]string{}
	err := GetUnreviewUser(&got)
	if err != nil {
		t.Errorf("获取未审核列表失败: %v", err)
		return
	}

	if len(got) != len(want) {
		t.Errorf("期待长度 %d, 实际长度 %d\n", len(want), len(got))
	}

	for _, v := range want {
		bool := false
		for _, w := range got {
			if v == w {
				bool = true
			}
		}
		if !bool {
			t.Errorf("结果中未包含该值 %s\n", v)
		}
	}
}

// 测试全部员工通过注册
func TestAllPass(t *testing.T) {
	err := AllPass()
	if err != nil {
		t.Errorf("通过失败: %v", err)
	}
}

// 测试部分员工通过注册
func TestPartPass(t *testing.T) {
	unreviewedList := [][2]string{
		{"13311674417", "qwer"},
		{"15146209355", "gerfrgv"},
	}
	err := PartPass(&unreviewedList, &[]int{0})
	if err != nil {
		t.Errorf("通过失败: %v", err)
	}
}

// 测试重置密码
func TestResetPwd(t *testing.T) {
	err := ResetPwd("13311674417", "admin")
	if err != nil {
		t.Errorf("通过失败: %v", err)
	}
}
