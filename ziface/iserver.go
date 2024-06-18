package ziface

// 定义服务器接口
type IServer interface {
	Start()
	Stop()
	Server()
	AddRouter(msgId uint32, router IRouter)  // 给当前服务注册一个路由业务方法，给客户端连接用
}