package ziface

type IServer interface {
	Start()
	Serve()
	Stop()
	// 路由功能：给当前的服务注册一个路由方法，供客户端连接处理使用
	AddRouter(msgid uint32, router IRouter)
	// 获取连接管理器
	GetConnMgr() IConnManager
	// 注册OnConnStart()的方法
	SetOnConnStart(func(IConnection))
	// 注册OnConnStop()的方法
	SetOnConnStop(func(IConnection))
	// 调用OnConnStart()的方法
	CallOnConnStart(IConnection)
	// 调用OnConnStop()的方法
	CallOnConnStop(IConnection)
}
