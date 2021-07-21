package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"fmt"
	"io"
	"os"
)

func main() {
	fileSlice := make(chan []byte, 64)

	go func() {
		f, _ := os.Open("a.zip")
		defer f.Close()

		stat, _ := f.Stat()
		fmt.Println("总字节：", stat.Size())

		if stat.Size() > 8388608 {
			buf := make([]byte, 8388608)
			for {
				n, err := f.Read(buf)
				if err == io.EOF {
					close(fileSlice)
					return
				} else if err != nil {
					fmt.Println("错误")
					return
				} else {
					fmt.Println("正常发送")
					fileSlice <- buf[:n]
				}
			}
		}
	}()

	for {
		slice, ok := <-fileSlice
		if !ok {
			return
		}
		sp := SP.NewSocketPacket(slice)
		fmt.Println(sp.String())
	}
}
