package main

import (
	"fmt"
	"zinx/znet"
	"zinx/ziface"
)

type PingRouter struct {
	znet.BaseRouter  // 嵌套后就可以继承了
}


func (pr *PingRouter) PreHandle(req ziface.IRequest){
	fmt.Println("Call Router PreHandle...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("Before ping...\n"))
	if err != nil {
		fmt.Println("Call back ping err...")
		return
	}
}

func (pr *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("Call back ping err...")
		return
	}
}

func (pr *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
	if err != nil {
		fmt.Println("Call back ping err...")
		return
	}
}

func main() {
	s := znet.NewServer("[Zinx v0.3]")

	s.AddRouter(&PingRouter{})

	s.Server()
}