package ziface

type IMsgHandle interface {
	DoMsgHandler(req IRequest)  // 以非阻塞方式处理消息， 调用Router中具体的接口
	AddRouter(msgId uint32, router IRouter)
}