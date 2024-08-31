package ziface

/*
IRequest 实现对已有连接和请求数据的封装
*/

type IRequest interface {
	// 得到当前连接
	GetConn() IConnection
	// 得到请求的数据
	GetMsgData() []byte
	// 得到请求msg的id
	GetMsgId() uint32
}
