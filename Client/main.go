package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Println("建立与服务端的链接")
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8888")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	conn.SetKeepAlive(true)
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}
	fmt.Println("连接成功")
	go recives(conn)

	txt := `网络不是法外之地，不是谣言滋生的温床，更不是某些人撒欢的自留田。
			7月16日凌晨，一网友在网上发布了杭州一名独居女孩遭遇“敲门杀人”事件的长文，因其描述的过程紧张，情节惨烈，一时间引起不少人转发讨论。16日下午，@杭州网警针对此事发布微博消息称，杭州出现敲门杀人案系谣言。
			从造谣到辟谣，不到1天的时间，围观者的情绪经历了过山车，还未从激愤、恐惧中走出来，又因被欺骗、利用而心生懊恼，多种情感混合交织在一起，可能比吃了苍蝇还恶心。
			不知道造谣者出于什么目的，要大半夜地编造这有鼻子有眼的谣言，更不知道造谣者有何居心，要如此费尽心思编文发图。有文字，有对话截图，很有迷惑性；有情节，有冲突，就更容易让人代入身份，投入情感。
			编造一条谣言，也许只是随意一想随手一发，但对社会的影响犹如蝴蝶效应，这起虚构的恶性事件，已经从个人社交平台传遍了网络。即便警方已经辟谣，其恶劣影响可能也很难即时消散。一些女性因此胆战心惊，心生畏惧，网上还出现了与谣言相匹配的“独居女性的安全神器”话题。大众的一本正经和造谣者的漫不经心形成了强烈的反差，让投入的人看起来像个“笑话”。
			也许造谣者看到有人被这则无中生有的谣言左右情绪时还沾沾自喜，或者为自己获得巨大的网络流量而引以为傲，但现实是，此时此刻，据@杭州网警消息，已经确认“是谣言了”。正如杭州网警所言，在网上造谣传谣，扰乱社会治安的行为需要承担法律责任。`

	for {
		fmt.Print("\n->")
		s := ""
		fmt.Scanln(&s)
		s = strings.TrimSpace(s)
		switch s {
		case "exit":
			return
		case "send":
			send(conn, []byte(txt))
		}
	}
}

func send(conn *net.TCPConn, b []byte) {
	b = append(b, []byte("\\ioEOF")...)
	buf := make([]byte, 32)

	breader := bytes.NewReader(b)
	//byteNum := 0
	for {
		n, err := breader.Read(buf)
		//byteNum += n
		if err == io.EOF {
			conn.Write(buf[:n])
			return
		}
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		// 直接将读到键盘输入数据，写到 socket 中，发送给服务器
		conn.Write(buf[:n])
	}
}

// 获取服务器回发数据
func recives(conn *net.TCPConn) {
	buf2 := make([]byte, 4096)
	for {
		n, err := conn.Read(buf2)
		if n == 0 {
			fmt.Println("服务器关闭了连接")
			return
		}
		if err != nil {
			fmt.Println("Read conn err:", err)
			return
		}
		fmt.Println("客户端读到：", string(buf2[:n]))
	}
}
