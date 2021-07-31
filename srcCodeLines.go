package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var lines int = 0

func main() {
	filepath.Walk(".", callback)
	fmt.Printf("项目中 *.go 代码 %d 行。\n", lines)
}

func callback(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	path = strings.ReplaceAll(path, "\\", "/")
	if info.IsDir() || filepath.Ext(info.Name()) != ".go" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	lines += countByte(b, '\n')

	return nil
}

func countByte(byteSlice []byte, b byte) int {
	if len(byteSlice) == 0 {
		return 0
	}
	count := 0

	for _, v := range byteSlice {
		if v == b {
			count++
		}
	}
	return count
}
