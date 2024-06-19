package ziface

type IMsgHandle interface {
	DoMsgHandler(req IRequest)  // 以非阻塞方式处理消息， 调用Router中具体的接口
	AddRouter(msgId uint32, router IRouter)
	StartWorkerPool()  // 启动worker工作池
	SendMsgToTaskQueue(req IRequest)  // 将消息交给task queue，由worker处理
}