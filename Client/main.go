package main

import (
	"encoding/json"
	"github.com/ying32/govcl/vcl"
	"io"
	"log"
	"os"
)

// 读取配置文件config.json，若文件不存在则主动创建；
// 配置文件可以手动设置tcp与udp服务器端口号；返工任务刷新时间间隔暂不生效
func init() {
	// 定义全局变量结构体
	config := struct {
		GlobalTCPPort        string `json:"远程TCP服务器端口"`
		GlobalUDPPort        string `json:"远程UDP服务器端口"`
		GlobalReworkInterval int    `json:"返工任务刷新时间间隔/小时"`
	}{
		GlobalTCPPort:        globalTCPPort,
		GlobalUDPPort:        globalUDPPort,
		GlobalReworkInterval: globalReworkInterval,
	}

	// 检查配置文件存在与否
	configFile := "config.json"
	stat, err := os.Stat(configFile)

	if os.IsNotExist(err) { // 配置文件不存在，创建配置文件
		f, err := os.Create(configFile)
		if err != nil {
			log.Printf("创建配置文件失败：%v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		b, _ := json.MarshalIndent(&config, "", "    ")
		_, err = f.Write(b)
		if err != nil {
			log.Println("写入配置文件失败")
			os.Exit(1)
		}
		return
	}

	// 配置文件存在
	if stat.IsDir() { // 配置文件名被占用
		log.Println("错误：文件名“config.json”被文件夹占用，请删除或重命名该文件夹后重试")
		os.Exit(1)
	}

	f, err := os.Open(configFile) // 打开配置文件
	if err != nil {
		log.Println("读取配置文件失败，使用默认参数")
		return
	}
	defer f.Close()

	b, err := io.ReadAll(f) // 读取配置文件
	if err != nil {
		log.Println("读取配置文件失败，使用默认参数")
		return
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Println("配置文件包含语法错误，使用默认参数")
		return
	}

	// 设置全局变量
	globalTCPPort = config.GlobalTCPPort
	globalUDPPort = config.GlobalUDPPort
	globalReworkInterval = config.GlobalReworkInterval
}

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	// vcl.Application.CreateForm(&MainForm)
	vcl.Application.CreateForm(&LoginForm)
	vcl.Application.CreateForm(&SignupForm)
	vcl.Application.SetShowMainForm(false)
	LoginForm.Show()

	go DirWatcher() // 监测项目文件夹

	vcl.Application.Run()
}
