package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	Version   string

	MaxPacketSize uint32 // 传输数据包的最大值
	MaxConn       int
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("zconf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, GlobalObject) // 将json数据解析到struct中
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "ZinxTest",
		Version:       "v0.4",
		TcpPort:       7778,
		Host:          "127.0.0.1",
		MaxConn:       1200,
		MaxPacketSize: 4096,
	}

	GlobalObject.Reload()
}
