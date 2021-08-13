package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var outputFile *os.File

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("将包名称用作参数")
		return
	}

	var err error
	outputFile, err = os.Create(args[1] + ".md")
	if err != nil {
		log.Println(err)
		return
	}
	defer outputFile.Close()

	err = filepath.Walk(".", walkFunc)
	if err != nil {
		log.Println(err)
	}
}

func walkFunc(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || filepath.Ext(info.Name()) != ".go" {
		return nil
	}

	err = genDoc(path)
	if err != nil {
		return err
	}

	return nil
}

func genDoc(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	t := fmt.Sprintf("## 文件：%s\r\n\r\n", path)
	printAndWrite(t)

	r := bufio.NewReader(f)

	//writeThrough := false
	bufList := []string{}
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s = strings.TrimSuffix(s, "\r\n")
		l := len(s)

		if l > 2 && s[:2] == "//" {
			bufList = append(bufList, strings.TrimSpace(s[2:]))
			continue
		}

		if l > 4 && s[:4] == "func" {
			printAndWrite("#### " + strings.TrimSuffix(s, " {") + "\r\n\r\n")
			printAndWrite(strings.Join(bufList, "\r\n\r\n") + "\r\n\r\n")
			bufList = []string{}
			continue
		}

	}
	return nil
}

func printAndWrite(s string) {
	fmt.Print(s)
	fmt.Fprint(outputFile, s)
}
