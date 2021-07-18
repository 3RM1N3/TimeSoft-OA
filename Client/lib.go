package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var globalName, globalID, globalIPAddr string

type Auth struct {
	Username string `json:"username"`
	Pwd      string `json:"password"`
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	ID   string `json:"id"`
}

type PrepareToUpload struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// 发送一个提交json数据的post请求登录账号
func postWithJsonLogin(addr, username, password string) (Response, error) {
	//post请求提交json数据
	auths := Auth{username, password}
	bytesJson, err := json.Marshal(auths)
	if err != nil {
		return Response{}, err
	}
	resp, err := http.Post(addr+"/login", "application/json", bytes.NewBuffer(bytesJson))
	if err != nil {
		return Response{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var result Response
	json.Unmarshal(body, &result)

	return result, nil
}

// 若输入的地址不含开头，补全开头的http://
func CompleteIPAddr(addr string) string {
	if !strings.HasPrefix(addr, "http://") ||
		!strings.HasPrefix(addr, "https://") {
		return "http://" + addr
	}
	return addr
}

// CheckMD5 计算md5值
func CheckMD5(s string) string {
	b := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(b))
}

// FolderChooser 打开一个文件夹选择窗口，返回选择的文件夹名字符串
func FolderChooser() string {
	var out bytes.Buffer

	cmd := exec.Command("powershell", "/c", "./folderChooser.bat")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}

	result := GbkToUtf8(out.Bytes())
	return strings.ReplaceAll(strings.TrimSpace(result), "\\", "/")
}

// GbcToUtf8 将gbk编码的字节切片转换为utf-8字符串
func GbkToUtf8(s []byte) string {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

// Zip 将文件或目录压缩为.zip文件
func Zip(srcFileOrDir string, destZip string) error {
	if _, err := os.Stat(destZip); err == nil { // 判断文件存在
		if err = os.Remove(destZip); err != nil {
			return err
		}
	}

	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	srcFileOrDir = strings.ReplaceAll(srcFileOrDir, "\\", "/")
	parentDir := path.Dir(srcFileOrDir) + "/"
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFileOrDir, func(everyFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		everyFilePath = strings.ReplaceAll(everyFilePath, "\\", "/")
		header.Name = strings.TrimPrefix(everyFilePath, parentDir)

		if info.IsDir() { // 如果是文件夹
			header.Name += "/"
			if _, err := archive.CreateHeader(header); err != nil {
				return err
			}
			return nil
		}

		// 如果是文件
		header.Method = zip.Deflate
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(everyFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return nil
}

/*
func get() {
	//get请求
	//http.Get的参数必须是带http://协议头的完整url,不然请求结果为空
	resp, _ := http.Get("http://localhost:8080/login2?username=admin&password=123456")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(body))
	fmt.Printf("Get request result: %s\n", string(body))
}

func postWithJson() {
	//post请求提交json数据
	auths := Auth{"admin", "123456"}
	ba, _ := json.Marshal(auths)
	resp, _ := http.Post("http://localhost:8080/login1", "application/json", bytes.NewBuffer([]byte(ba)))
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Post request with json result: %s\n", string(body))
}

func postWithUrlencoded() {
	//post请求提交application/x-www-form-urlencoded数据
	form := make(url.Values)
	form.Set("username", "admin")
	form.Add("password", "123456")
	resp, _ := http.PostForm("http://localhost:8080/login2", form)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Post request with application/x-www-form-urlencoded result: %s\n", string(body))
}
*/
