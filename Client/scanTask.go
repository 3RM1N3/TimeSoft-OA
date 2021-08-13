package main

import (
	"TimeSoft-OA/lib"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"
)

var ErrScanTooFast = errors.New("你的扫描速度似乎有些快于常人，为判定是否作弊，请主动与管理员取得联系。")

// 获取当前已有客户
func GetClientCo() ([]string, error) {
	b, err := sendUDPMsg(globalServerAddr, lib.ClientCo, nil)
	if err != nil {
		return nil, err
	}

	if len(b) == 1 {
		return nil, lib.ReportCode(b[0]).ToError()
	}

	clientCos := []string{}
	err = json.Unmarshal(b, &clientCos)
	if err != nil {
		return nil, err
	}

	return clientCos, nil
}

// 设置项目文件夹
func SetProjectDir(dirPath string) error {
	dirEntryList, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("读取目标文件夹失败")
		return err
	}
	if len(dirEntryList) == 0 { // 空文件夹能够直接使用
		return nil
	}

	// 验证.verf文件
	if r, err := CheckVerf(dirPath); err != nil {
		return err
	} else if !r {
		return errors.New("文件夹校验失败，请选择曾用的项目文件夹或空文件夹")
	}

	return nil
}

// 校验.verf文件
func CheckVerf(dirPath string) (bool, error) {
	verfPath := path.Join(dirPath, ".verf")
	f, err := os.Open(verfPath)
	if err != nil {
		log.Println("校验文件打开失败")
		return false, err
	}

	b := make([]byte, 40)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	wantVerfStr := string(b)
	gotVerfStr, err := GenVerf(dirPath, -1)
	if err != nil {
		log.Println("生成校验码失败")
		return false, err
	}
	if gotVerfStr != wantVerfStr {
		return false, nil
	}

	return true, nil
}

// 生成.verf校验字符串
// fileNum设置为-1，则读取文件夹生成字符串，否则使用fileNum值
func GenVerf(dirPath string, fileNum int) (string, error) {
	if fileNum == -1 {
		fileNum = 0
		err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Name() == ".verf" {
				return nil
			} else if !info.IsDir() && filepath.Ext(info.Name()) != ".jpg" {
				return nil
			}

			fileNum++
			return nil
		})

		if err != nil {
			return "", err
		}
	}

	macs, err := GetMacAddrs()
	if err != nil {
		return "", err
	}
	sort.Strings(macs)

	got := fmt.Sprintf("%s|%s|%d", "!g*657#JW@$", macs[0], fileNum)

	return "verify!#" + lib.MD5(got), nil
}

// 监测项目文件夹，increaseCh可以传出增长的文件数，errCh传出错误，
// scanOver设置为true时扫描结束
func DirWatcher(dirPath string, increaseCh chan int, errCh chan error, scanOver *bool) {

	firstWatch := true
	existNotJpg := false
	preNum, fileNum := 0, 0

	defer func() {
		// 生成校验文件
		s, err := GenVerf(dirPath, preNum)
		if err != nil {
			log.Println("生成校验字符串失败，请与管理员取得联系。")
			errCh <- err
			return
		}
		f, err := os.OpenFile(filepath.Join(dirPath, ".verf"), os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Println("打开校验文件失败，请与管理员取得联系。")
			errCh <- err
			return
		}
		_, err = f.Write([]byte(s))
		if err != nil {
			log.Println("写入校验文件失败，请与管理员取得联系。")
			errCh <- err
		}
	}()

	for !*scanOver {
		time.Sleep(5 * time.Second)

		existNotJpg = false
		err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == ".verf" {
				return nil
			}
			if !info.IsDir() && filepath.Ext(info.Name()) != ".jpg" {
				existNotJpg = true
				return nil
			}
			fileNum++
			return nil
		})
		if err != nil {
			errCh <- err
			return
		}

		if existNotJpg {
			log.Println("检测到非*.jpg格式文件，请确认扫描设置正确，删除格式错误的文件后重试。")
		}

		if firstWatch { // 第一次读取不判断速度
			firstWatch = false
			preNum, fileNum = fileNum, 0
			continue
		}

		d := fileNum - preNum
		if d > 7 {
			log.Println("文件生成速度过快")
			errCh <- ErrScanTooFast
			return
		}
		if d >= 0 {
			increaseCh <- d
		}

		preNum, fileNum = fileNum, 0
	}
}

// 扫描结束，打包文件夹并提交，dirPath为要提交的文件夹，scanOver设置结束监测文件夹
func ScanOverPackSubmit(dirPath, clientCo, uploader string, scanOrEdit byte, scanOver *bool) error {
	*scanOver = true // 结束监测

	tempFile := ".temp_archive" + uploader
	err := lib.Zip(dirPath, tempFile)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", globalServerAddr+globalTCPPort)
	if err != nil {
		return err
	}
	fmt.Println("与服务器连接成功，发送文件")

	err = lib.SendFile(tempFile, clientCo, uploader, scanOrEdit, conn)
	if err != nil {
		return err
	}

	return nil
}

// 获取今日工作量
func TodayWorkload() (lib.WorkLoadJson, error) {
	var todayWorkload = lib.WorkLoadJson{
		Phone: globalPhone,
	}
	b, err := sendUDPMsg(globalServerAddr, lib.WorkLoad, todayWorkload)
	if err != nil {
		return todayWorkload, err
	}

	if len(b) == 1 {
		return todayWorkload, lib.ReportCode(b[0]).ToError()
	}

	err = json.Unmarshal(b, &todayWorkload)
	if err != nil {
		return todayWorkload, err
	}

	return todayWorkload, nil
}
