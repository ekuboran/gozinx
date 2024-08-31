package znet

import "gozinx/ziface"

type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

// 得到当前连接
func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

// 得到请求的数据
func (r *Request) GetMsgData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
