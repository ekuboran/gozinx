package znet

/*
连接模块
*/

import (
	"errors"
	"gozinx/utils"
	"gozinx/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TCPServer    ziface.IServer
	Conn         *net.TCPConn
	ConnID       uint32
	connState    bool
	ExitChan     chan bool
	msgChan      chan []byte
	MsgHandler   ziface.IMsgHandler     // 消息管理模块，绑定msgId和对应的路由业务
	property     map[string]interface{} // 连接属性集合
	propertyLock sync.RWMutex
}

// 初始化连接模块
func NewConnection(server ziface.IServer, conn *net.TCPConn, connid uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connid,
		connState:  true,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	// 将conn加入当当前server的Connmanager中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	utils.DPrintln("Reader goroutine is running")
	defer utils.DPrintln("ConnID: ", c.ConnID, " Reader has exited, RemoterAddr: ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 创建拆包封包对象
		dp := NewDataPack()
		// 读头信息共8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			utils.DPrintln("server read head err: ", err)
			break
		}

		// 拆包获得msgID和msgLen
		msg, err := dp.Unpack(headData)
		if err != nil {
			utils.DPrintln("server unpack err: ", err)
			break
		}

		// 根据msgLen读取msgData
		var data []byte
		if msg.GetDataLen() > 0 {
			// 如果有数据，进行第二次读取
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				utils.DPrintln("server read data err: ", err)
				break
			}
		}
		msg.SetData(data)

		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 将消息发送给消息队列处理
		c.MsgHandler.SendMsgToTaskQueue(req)
	}

}

func (c *Connection) StartWriter() {
	utils.DPrintln("Writer goroutine is running")
	defer utils.DPrintln("ConnID: ", c.ConnID, " Writer has exited, RemoterAddr: ", c.GetRemoteAddr().String())
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				utils.DPrintln("Send data error")
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 启动连接
func (c *Connection) Start() {
	utils.DPrintln("Conn Start, ConnID: ", c.ConnID)

	// 启动从当前连接读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 执行开发者自定义的hook函数——创建连接之后调用的处理业务
	c.TCPServer.CallOnConnStart(c)
}

// 停止连接
func (c *Connection) Stop() {
	utils.DPrintln("Conn Stop, ConnID: ", c.ConnID)
	// 如果连接已经断开
	if !c.connState {
		return
	}
	c.connState = false

	// 执行开发者自定义的hook函数——销毁连接之前调用的处理业务
	c.TCPServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	// 将当前连接从Connmanager中移除
	c.TCPServer.GetConnMgr().Delete(c)

	close(c.ExitChan)
	close(c.msgChan)
}

// 获取连接句柄TCP socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取连接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

// 获取连接对应的远端客户端ip:port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if !c.connState {
		return errors.New("Connection has closed")
	}
	dp := NewDataPack()
	msg := NewMessage(msgId, data)
	sendData, err := dp.Pack(msg)
	if err != nil {
		utils.DPrintln("server send data pack err:", err)
		return errors.New("pack error msg")
	}

	// 将数据发给Writer
	c.msgChan <- sendData
	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property FOUND")
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
