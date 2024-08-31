package znet

/*
Server模块
*/

import (
	"fmt"
	"gozinx/utils"
	"gozinx/ziface"
	"net"
)

type Server struct {
	Name        string
	IPVersion   string
	Ip          string
	Port        int
	MsgHandler  ziface.IMsgHandler       // 当前server的消息管理模块，绑定msgId和对应的路由业务
	Connmanager ziface.IConnManager      // 当前server的连接管理模块
	OnConnStart func(ziface.IConnection) // Server创建连接之后自动调用的Hook函数—OnConnStart()
	OnConnStop  func(ziface.IConnection) // Server销毁连接之前自动调用的Hook函数—OnConnStop()
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:        utils.GloablObject.Name,
		IPVersion:   "tcp4",
		Ip:          utils.GloablObject.Host,
		Port:        utils.GloablObject.TCPPort,
		MsgHandler:  NewMsgHandler(),
		Connmanager: NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	utils.DPrintf("[Zinx] Server name: %s listening at %s:%d is starting\n", s.Name, s.Ip, s.Port)
	utils.DPrintf("[Zinx] Version: %s, MaxConn: %d, Maxpackagesize: %d\n", utils.GloablObject.Version, utils.GloablObject.MaxCon, utils.GloablObject.MaxPackageSize)

	// 开启工作池
	s.MsgHandler.StartWorkerPool()

	go func() {
		// 获取一个tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			utils.DPrintln("resolve tcp addr error:", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			utils.DPrintln("listen ", s.IPVersion, "err: ", err)
			return
		}

		utils.DPrintln("start zinx server success", s.Name, "success, listening")

		var cid uint32

		// 阻塞地等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				utils.DPrintln("Accept err: ", err)
				continue
			}

			// 判断是否超出最大连接数量
			if s.Connmanager.GetLen() >= utils.GloablObject.MaxCon {
				// TODO 给客户端响应一个超出最大连接数量的错误信息
				utils.DPrintln("Too many connections")
				conn.Close()
				continue
			}

			// 将处理新链接的业务方法和 conn绑定得到 我们定义的连接模块，并将该连接加入到Connmanager中
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Serve() {
	s.Start()

	// TODO 一些扩展业务

	select {}
}

func (s *Server) Stop() {
	// 将服务器的资源、状态等停止或回收
	utils.DPrintf("[Zinx] Server name: %s Stoped ", s.Name)
	s.Connmanager.ClearConn()
}

func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgid, router)
	utils.DPrintln("Router add success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.Connmanager
}

// 注册OnConnStart()的方法
func (s *Server) SetOnConnStart(hookfunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookfunc
}

// 注册OnConnStop()的方法
func (s *Server) SetOnConnStop(hookfunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookfunc
}

// 调用OnConnStart()的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		utils.DPrintln("———>Call OnConnStart()")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop()的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		utils.DPrintln("———>Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
