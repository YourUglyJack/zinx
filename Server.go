package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// 自定义路由
type PingRouter struct {
	znet.BaseRouter  // 嵌套后就可以继承了
}


// func (pr *PingRouter) PreHandle(req ziface.IRequest){
// 	fmt.Println("Call Router PreHandle...")
// 	_, err := req.GetConnection().GetTCPConnection().Write([]byte("Before ping...\n"))
// 	if err != nil {
// 		fmt.Println("Call back ping err...")
// 		return
// 	}
// }

func (pr *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	
	// 读取客户端的数据
	fmt.Println("recv from client, msg Id:", req.GetMsgId(), "data:", string(req.GetData()))
	
	// 回写数据
	err := req.GetConnection().SendMsg(1, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
	
	// _, err := req.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	// if err != nil {
	// 	fmt.Println("Call back ping err...")
	// 	return
	// }
}

// func (pr *PingRouter) PostHandle(req ziface.IRequest) {
// 	fmt.Println("Call Router Handle...")
// 	_, err := req.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
// 	if err != nil {
// 		fmt.Println("Call back ping err...")
// 		return
// 	}
// }

func main() {
	s := znet.NewServer("[Zinx v0.5]")

	s.AddRouter(&PingRouter{})

	s.Server()
}