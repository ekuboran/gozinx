package znet

import "gozinx/ziface"

// 实现Router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法重写 (接口隔离)
type BaseRouter struct {
}

// 处理业务之前的钩子方法Hook
func (r *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理业务时的主法Hook
func (r *BaseRouter) Handle(request ziface.IRequest) {}

// 处理业务之后的的钩子方法Hook
func (r *BaseRouter) PostHandle(request ziface.IRequest) {}
