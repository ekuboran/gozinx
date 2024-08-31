package utils

import (
	"encoding/json"
	"gozinx/ziface"
	"os"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些用户自选的参数写入zinx.json
*/

type GlobalObj struct {
	// Server
	TCPServer ziface.IServer
	Name      string
	Host      string
	TCPPort   int

	// Zinx
	Version        string
	MaxCon         int
	MaxPackageSize uint32
	WorkerPoolSize uint32 // worker池大小
	MaxTaskLen     uint32 // 消息队列的最大长度
}

var GloablObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到对象中
	err = json.Unmarshal(data, &GloablObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 没有设置配置文件时的默认参数
	GloablObject = &GlobalObj{
		Name:           "ZinxServer",
		Host:           "0.0.0.0",
		TCPPort:        8999,
		Version:        "V0.4",
		MaxCon:         10,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxTaskLen:     1024,
	}

	// 从配置文件中加载配置参数
	GloablObject.Reload()
}
