package lib

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Uint16ToByte(i uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToUint16(b []byte) (uint16, error) {
	if len(b) == 0 {
		return 0, nil
	}
	i := uint16(0)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func Uint32ToByte(i uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToUint32(b []byte) (uint32, error) {
	if len(b) == 0 {
		return 0, nil
	}
	i := uint32(0)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Zip 将目录下的内容压缩为.zip文件
func Zip(srcDir, destZip string) error {
	srcDir = filepath.Clean(srcDir)

	if _, err := os.Stat(destZip); err == nil { // 判断文件存在
		if err = os.Remove(destZip); err != nil { // 不存在则移除
			return err
		}
	}

	zipfile, err := os.Create(destZip) // 创建压缩文件
	if err != nil {
		return err
	}
	defer zipfile.Close() // 函数返回后关闭文件

	srcDir = strings.ReplaceAll(srcDir, "\\", "/") // 将windows路径中的反斜杠替换成斜杠
	parentDir := srcDir + "/"                      // 获取全部上级文件夹
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = filepath.Walk(srcDir, func(everyFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && everyFilePath == srcDir {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		everyFilePath = strings.ReplaceAll(everyFilePath, "\\", "/")
		header.Name = strings.TrimPrefix(everyFilePath, parentDir)
		log.Println(header.Name)

		if info.IsDir() { // 如果是文件夹
			// header.Name += "/"
			// if _, err := archive.CreateHeader(header); err != nil { // 创建header
			// 	return err
			// }
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
	if err != nil {
		return err
	}

	return nil
}

// Unzip 将.zip文件解压至目录，如果目录不存在则自动创建
func Unzip(srcZip, destDir string) ([]string, error) {
	destDir = path.Clean(destDir) // 整理目录字符串
	fileIDList := []string{}      // 初始化档号列表

	r, err := zip.OpenReader(srcZip) // 打开源文件
	if err != nil {
		return nil, err
	}
	defer r.Close() // 函数返回后关闭源文件

	// 逐个枚举压缩文档内的文件
	for _, f := range r.File {
		pathInArchive, err := GbkToUtf8([]byte(f.Name)) // 获取压缩包内的文件路径
		if err != nil {
			return nil, err
		}

		fullPath := path.Join(destDir, pathInArchive) // 生成完整路径
		ParentDir := path.Dir(fullPath)               // 获取上级目录

		if path.Dir(ParentDir) == destDir {
			fileIDList = append(fileIDList, path.Base(ParentDir)) // 插入档号列表
		}

		if ParentDir != "." {
			err := os.MkdirAll(ParentDir, 0666) // 创建全部上级目录
			if err != nil {
				return nil, err
			}
		}
		// fmt.Printf("得到 %s\n", pathInArchive)
		rc, err := f.Open()
		if err != nil {
			log.Printf("包内文件%s读取失败 %v\n", pathInArchive, err)
			continue
		}
		destFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0755) // 创建目标文件
		if err != nil {
			log.Printf("创建本地文件%s失败 %v\n", fullPath, err)
			continue
		}
		_, err = io.Copy(destFile, rc) // 复制文件
		if err != nil {
			log.Printf("包内文件%s解压失败 %v\n", pathInArchive, err)
			continue
		}
		destFile.Close()
		rc.Close()
	}
	return fileIDList, nil
}

// 将gbk编码的字符串转码为utf-8
func GbkToUtf8(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

// 计算字符串的md5值
func MD5(s string) string {
	b := []byte(s)
	return fmt.Sprintf("%X", md5.Sum(b))
}
