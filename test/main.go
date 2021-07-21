package main

import (
	SP "TimeSoft-OA/SocketPacket"
	"fmt"
	"os"
)

func main() {
	fileSliceChan := make(chan SP.SocketPacket, 64)

	go func() {
		f, _ := os.Open("a.zip")
		defer f.Close()

		SP.NewZipPacket(f, fileSliceChan)
	}()

	for {
		sp := <-fileSliceChan
		fmt.Print(sp.String(), "\n\n")
	}
}
