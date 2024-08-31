package znet

import (
	"errors"
	"gozinx/utils"
	"gozinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 新增连接
	cm.connections[conn.GetConnId()] = conn

	utils.DPrintln("connID=", conn.GetConnId(), " Add success, conn num=", cm.GetLen())
}

// 删除连接
func (cm *ConnManager) Delete(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除连接信息
	delete(cm.connections, conn.GetConnId())
	utils.DPrintln("connID=", conn.GetConnId(), " Remove from ConnManager success, conn num=", cm.GetLen())
}

// 根据ConnID获取连接
func (cm *ConnManager) GetConn(connid uint32) (ziface.IConnection, error) {
	// 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connid]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not FOUND")
}

// 获取连接总个数
func (cm *ConnManager) GetLen() int {
	return len(cm.connections)
}

// 清除全部连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connid, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connid)
	}

	utils.DPrintln("Clear all conn success! conn num=", cm.GetLen())
}
