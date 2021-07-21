package main

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	globalName, globalID, globalIPAddr string
)

type Auth struct {
	Username string `json:"username"`
	Pwd      string `json:"password"`
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	ID   string `json:"id"`
}

type ScanedJob struct {
	JobID        string `json:"jobid"`
	FolderName   string `json:"foldername"`
	SubFolderNum int    `json:"subfoldernum"`
	AllFileNum   int    `json:"allfilenum"`
	JobType      string `json:"jobtype"`
	UploadTime   int    `json:"uploadtime"`
}

// CheckMD5 计算md5值
func CheckMD5(s string) string {
	b := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(b))
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
