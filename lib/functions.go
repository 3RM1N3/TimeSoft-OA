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

// Zip 将文件或目录压缩为.zip文件
func Zip(srcFileOrDir, destZip string) error {
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

// Unzip 将.zip文件解压至目录，如果目录不存在则自动创建
func Unzip(srcZip, destDir string) error {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(destDir, 0644)
	if err != nil {
		log.Println("创建目标文件夹失败", err)
		return err
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		pathInArchive, _ := GbkToUtf8([]byte(f.Name))
		fullPath := path.Join(destDir, pathInArchive)
		dir := path.Dir(pathInArchive)
		if dir != "." {
			err := os.MkdirAll(path.Join(destDir, dir), 0644)
			if err != nil {
				log.Println("创建包内文件夹结构失败", err)
				continue
			}
		}
		fmt.Printf("Contents of %s:\n", pathInArchive)
		rc, err := f.Open()
		if err != nil {
			log.Printf("包内文件%s读取失败 %v\n", pathInArchive, err)
			continue
		}
		destFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Printf("创建本地文件%s失败 %v\n", fullPath, err)
			continue
		}
		_, err = io.Copy(destFile, rc)
		if err != nil {
			log.Printf("包内文件%s解压失败 %v\n", pathInArchive, err)
			continue
		}
		destFile.Close()
		rc.Close()
		fmt.Print("\n\n")
	}
	return nil
}

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
