package ziface

import "net"

type IConnection interface {
	// 启动连接
	Start()

	// 停止连接
	Stop()

	// 获取连接句柄TCP socket
	GetTCPConnection() *net.TCPConn

	// 获取连接ID
	GetConnId() uint32

	// 获取连接对应的远端客户端ip:port
	GetRemoteAddr() net.Addr

	// 发送数据
	SendMsg(msgId uint32, data []byte) error

	// 设置连接属性
	SetProperty(string, interface{})

	// 获取连接属性
	GetProperty(string) (interface{}, error)

	// 移除连接属性
	RemoveProperty(string)
}

// 定义一个连接所绑定业务的处理函数类型，参数为（处理内容，处理长度）
type HandleFunc func(*net.TCPConn, []byte, int) error
