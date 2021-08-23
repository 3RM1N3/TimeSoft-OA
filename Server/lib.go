package main

import (
	"TimeSoft-OA/lib"
	"archive/zip"
	"database/sql"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ProjectPath = "."
	PortUDP = ":8080"
	PortTCP = ":8888"
)
var db *sql.DB
var UserEditMap = map[string][]string{}

// ZipDirs 压缩多个目录进一个文件
func ZipDirs(dirList []string, destZip string) error {
	log.Println("压缩文件", destZip)
	if _, err := os.Stat(destZip); err == nil { // 判断文件存在
		if err = os.Remove(destZip); err != nil { // 存在则移除
			return err
		}
	}

	zipfile, err := os.Create(destZip) // 创建压缩文件
	if err != nil {
		return err
	}
	defer zipfile.Close() // 函数返回后关闭文件
	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	for _, srcDir := range dirList {
		srcDir = filepath.Join(ProjectPath, srcDir)

		srcDir = strings.ReplaceAll(srcDir, "\\", "/") // 将windows路径中的反斜杠替换成斜杠
		parentDir := filepath.Dir(srcDir) + "/"        // 获取全部上级文件夹

		err = filepath.Walk(srcDir, func(everyFilePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() { // 如果是文件夹
				return nil
			}

			if info.Name() == ".verf" {
				return nil
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			everyFilePath = strings.ReplaceAll(everyFilePath, "\\", "/") // 实际文件路径
			header.Name = strings.TrimPrefix(everyFilePath, parentDir)   // 压缩包内路径
			log.Println(header.Name)

			// 如果是文件
			header.Method = zip.Deflate
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(everyFilePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			file.Close()
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUnreviewUser 获取未审核员工信息
func GetUnreviewUser(unreviewedUserList *[][2]string) error {
	rows, err := db.Query(`SELECT PHONE, REALNAME FROM UNREVIEWED_USER`)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Println("查询数据库失败")
		return err
	}

	for rows.Next() {
		tempList := [2]string{}
		err = rows.Scan(&tempList[0], &tempList[1])
		if err != nil {
			log.Println("获取某条内容失败")
			return err
		}
		*unreviewedUserList = append(*unreviewedUserList, tempList)
	}

	return nil
}

// AllPass 令未审核的员工全部通过
func AllPass() error {
	_, err := db.Exec(`INSERT INTO USER SELECT * FROM UNREVIEWED_USER`)
	if err != nil {
		log.Println("未审核员工通过失败")
		return err
	}
	_, err = db.Exec(`DELETE FROM UNREVIEWED_USER`)
	if err != nil {
		log.Println("清空未审核表失败，请手动清空")
		return err
	}

	return nil
}

// PartPass 部分通过
func PartPass(unreviewedList *[][2]string, passedIndex *[]int) error {
	for i, v := range *unreviewedList {
		if intListContains(*passedIndex, i) {
			// 通过用户插入至用户表
			_, err := db.Exec(`INSERT INTO USER SELECT * FROM UNREVIEWED_USER WHERE PHONE = ?`, v[0])
			if err != nil {
				log.Println("插入失败")
				return err
			}
		}
		// 删除已审核记录
		_, err := db.Exec(`DELETE FROM UNREVIEWED_USER WHERE PHONE = ?`, v[0])
		if err != nil {
			log.Println("UNREVIEWED_USER中删除此条记录失败，请手动删除 电话号码:", v[0])
			return err
		}
	}
	return nil
}

// 整型切片中是否包含某个数字
func intListContains(l []int, i int) bool {
	for _, e := range l {
		if e == i {
			return true
		}
	}
	return false
}

// ResetPwd 重置密码
func ResetPwd(phone, newPwd string) error {
	md5pwd := lib.MD5(newPwd)

	_, err := db.Exec(`UPDATE USER SET PWD = ? WHERE PHONE = ?`, md5pwd, phone)
	if err != nil {
		log.Println("重置失败")
		return err
	}

	return nil
}

// 获取文件夹内文件数
func getSubFileNum(path string) (int, error) {
	num := 0

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		num++

		return nil
	})

	return num, err
}
