package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32] ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
		panic("repeated api, msgId:" +  strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("add api msgId:", msgId)
}