package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	//a()
	//b()
	c()
}

func c() {
	b1 := make([]byte, 5)
	fmt.Println(len(b1))
	//SP.TrimByteRight0x0(&b1)
	for _, b := range b1 {
		fmt.Printf("%X\n", b)
	}

}

// 测试解压文件
func b() {
	fmt.Println(path.Dir("a.txt"))
	err := SP.Unzip("a48.6.zip", "./folder")
	if err != nil {
		log.Println("解压失败", err)
	}
}

// 测试传输压缩文件
func a() {
	fileSliceChan := make(chan SP.SocketPacket, 512)

	f, _ := os.Open("b.zip")
	defer f.Close()

	SP.NewZipPacket(f, fileSliceChan)

	for i := 0; i < 10; i++ {
		fmt.Printf("%d\n", 10-i)
		time.Sleep(time.Second)
	}

	sp := <-fileSliceChan
	if sp.TypeByte == SP.FileUpload {
		sp.SplicingFile(fileSliceChan)
	}
	fmt.Println("Done.")
}
