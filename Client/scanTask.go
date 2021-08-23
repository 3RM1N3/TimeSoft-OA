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
	"time"

	"github.com/ying32/govcl/vcl"
)

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
		return errors.New("校验码不匹配")
	}

	return nil
}

// 校验.verf文件
func CheckVerf(dirPath string) (bool, error) {
	verfPath := path.Join(dirPath, ".verf")
	f, err := os.Open(verfPath)
	if err != nil {
		return false, errors.New("校验文件不存在")
	}
	defer f.Close()

	b := make([]byte, 50)
	n, err := f.Read(b)
	if err != nil {
		return false, err
	}
	wantVerfStr := string(b[:n])
	gotVerfStr, err := GenVerf(dirPath, -1)
	if err != nil {
		return false, errors.New("生成校验码失败")
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

			if path == dirPath {
				return nil
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

	got := fmt.Sprintf("%s|%s|%d", "!g*657#JW@$", macAddr, fileNum)

	return "verify!#" + lib.MD5(got), nil
}

// 循环监测项目文件夹
func DirWatcher() {
	for {
		_ = <-ChStartScan
		preNum := 0
		OverScan = false
		d := 0
		for !OverScan {
			i, err := WatchDir(ProjectDir)
			if err != nil {
				vcl.ThreadSync(func() {
					vcl.ShowMessageFmt("%v", err)
					MainForm.StatusBar.SetSimpleText("就绪")
				})
				ProjectDir = ""
				break
			}
			if preNum == 0 {
				preNum = i
				continue
			}
			d = i - preNum
			if d > 7 {
				os.Remove(filepath.Join(ProjectDir, ".verf"))
				vcl.ThreadSync(func() {
					vcl.ShowMessageFmt("%v", ErrScanTooFast)
					MainForm.StatusBar.SetSimpleText("就绪")
				})
				ProjectDir = ""
				break
			}
			vcl.ThreadSync(func() {
				t := fmt.Sprintf("当前速度%d个/分，现有%d个文件（夹）", d*12, i)
				MainForm.StatusBar.SetSimpleText(t)
			})
			preNum = i
			time.Sleep(5 * time.Second)
		}
	}
}

// 传入监测项目文件夹，生成校验文件并返回项目内文件数量或错误信息
func WatchDir(dirPath string) (int, error) {
	fileNum := 0

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dirPath {
			return nil
		}
		if info.Name() == ".verf" {
			return nil
		}
		if !info.IsDir() && filepath.Ext(info.Name()) != ".jpg" {
			return ErrFindNotJpg
		}

		fileNum++
		return nil
	})
	if err != nil {
		return 0, err
	}

	verf, _ := GenVerf(dirPath, fileNum)

	f, err := os.OpenFile(filepath.Join(dirPath, ".verf"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	f.WriteString(verf)

	return fileNum, nil
}

// 扫描结束，打包文件夹并提交，dirPath为要提交的文件夹，scanOver设置结束监测文件夹
// 0为扫描1修图
func PackSubmitDir(dirPath, clientCo, uploader string, scanOrEdit byte) error {

	tempFile := ".temp_archive" + uploader // 临时文件名
	defer os.Remove(tempFile)              // 函数返回后删除临时文件
	err := lib.Zip(dirPath, tempFile)      // 压缩文件夹
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", globalServerAddr+globalTCPPort) // 建立TCP连接
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

// 下载并解压修图任务
func SaveAndUnzip(dirPath string) ([]string, error) {
	conn, err := net.Dial("tcp", globalServerAddr+globalTCPPort)
	if err != nil {
		return nil, err
	}
	fmt.Println("与服务器连接成功，接收文件")
	name, err := ClientReceiveFile([]string{}, conn)
	if err != nil {
		return nil, err
	}
	defer os.Remove(name)
	fmt.Println("收到文件", name)
	fileIDs, err := lib.Unzip(name, dirPath)
	if err != nil {
		return nil, err
	}
	return fileIDs, nil
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
