package ziface

/*
  路由抽象接口
  路由的数据就是Request
*/

type IRouter interface {
	// 处理业务之前的钩子方法Hook
	PreHandle(request IRequest)

	// 处理业务时的主法Hook
	Handle(request IRequest)

	// 处理业务之后的的钩子方法Hook
	PostHandle(request IRequest)
}
