package ziface

type IConnManager interface {
	// 添加连接
	Add(IConnection)
	// 删除连接
	Delete(IConnection)
	// 根据ConnID查询连接
	GetConn(uint32) (IConnection, error)
	// 获取连接总个数
	GetLen() int
	// 清除全部连接
	ClearConn()
}
