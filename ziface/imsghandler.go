package ziface

type IMsgHandler interface {
	// 执行对应Router的消息处理方法
	DoMsgHandle(request IRequest)
	// 添加新的具体的处理逻辑
	AddRouter(msgid uint32, router IRouter)
	// 启动Worker工作池
	StartWorkerPool()
	// 将消息发送给TaskQueue
	SendMsgToTaskQueue(IRequest)
}
