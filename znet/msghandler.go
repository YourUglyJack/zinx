package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize), // 一个worker对应一个队列
	}
}

func (mh *MsgHandle) DoMsgHandler(req ziface.IRequest) {
	handler, ok := mh.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("api msgId:", req.GetMsgId(), "is not found")
		return
	}

	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgId:" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("add api msgId:", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker Id:", workerId, "is start")
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	workerId := req.GetConnection().GetConnId() % mh.WorkerPoolSize
	fmt.Println("add connId", req.GetConnection().GetConnId(), "req msgId", req.GetMsgId(), "to workerId", workerId)
	mh.TaskQueue[workerId] <- req // 将请求发给任务队列
}
