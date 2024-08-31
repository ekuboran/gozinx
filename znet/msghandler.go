package znet

import (
	"gozinx/utils"
	"gozinx/ziface"
	"strconv"
)

type MsgHandler struct {
	ApiM           map[uint32]ziface.IRouter
	TaskQueue      []chan ziface.IRequest // 消息队列
	WorkerPoolSize uint32                 // worker池大小
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		ApiM:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GloablObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GloablObject.WorkerPoolSize),
	}
}

// 执行对应Router的消息处理方法
func (mh *MsgHandler) DoMsgHandle(request ziface.IRequest) {
	handler, ok := mh.ApiM[request.GetMsgId()]
	if !ok {
		utils.DPrintln("api msgId=", request.GetMsgId(), "is NOT FOUND! NEED REGISTER")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 添加新的具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgid uint32, router ziface.IRouter) {
	if _, ok := mh.ApiM[msgid]; ok {
		panic("repeat api, msgId=" + strconv.Itoa(int(msgid)))
	}
	mh.ApiM[msgid] = router
}

// 启动一个worker工作池，每个zinx框架只启动一次
func (mh *MsgHandler) StartWorkerPool() {
	// 根据worker池的大小开辟消息任务队列
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 每个worker对应一个消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GloablObject.MaxTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个工作流程
func (mh *MsgHandler) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	utils.DPrintln("WorkerID = ", workerID, "is started...")
	for request := range taskQueue {
		mh.DoMsgHandle(request)
	}
}

// 将消息发送给TaskQueue
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1.将消息平均分配给不同的worker
	// 根据客户端建立的ConnID分配（更正确的应该是根据request id分配，这里是因为没设置requestid）
	workerId := request.GetConn().GetConnId() % mh.WorkerPoolSize
	utils.DPrintln("Add CoonID=", request.GetConn().GetConnId(), " request MsgID=", request.GetMsgId(), " to WorkerID=", workerId)

	// 2.将消息放松给对应worker的TaskQueue
	mh.TaskQueue[workerId] <- request
}
