package main

import (
	"TimeSoft-OA/lib"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// 获取未审核员工信息
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

// 令未审核的员工全部通过
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

// 部分通过
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

func intListContains(l []int, i int) bool {
	for _, e := range l {
		if e == i {
			return true
		}
	}
	return false
}

// 重置密码
func ResetPwd(phone, newPwd string) error {
	md5pwd := lib.MD5(newPwd)

	_, err := db.Exec(`UPDATE USER SET PWD = ? WHERE PHONE = ?`, md5pwd, phone)
	if err != nil {
		log.Println("重置失败")
		return err
	}

	return nil
}
