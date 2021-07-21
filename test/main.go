package main

import (
	"fmt"
	"os"
)

func main() {
	//fileSlice := make(chan []byte, 64)

	f, _ := os.Open("a.zip")
	stat, _ := f.Stat()
	defer f.Close()
	var partsNum int = 1/8388608 + 1
	fmt.Println(partsNum, stat.Size())
	/*
		go func() {
			f, _ := os.Open("a.zip")
			defer f.Close()

			newPacket(f, fileSlice)

			stat, _ := f.Stat()
			fmt.Println("总字节：", stat.Size())

			if stat.Size() > 8388608 {
				buf := make([]byte, 8388608)
				for {
					n, err := f.Read(buf)
					if err != nil {
						close(fileSlice)
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
	*/
}

func newPacket(f *os.File, ch chan []byte) {

}
