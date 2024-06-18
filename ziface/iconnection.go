package ziface

import "net"

// 定义连接接口
type IConnection interface {
	Start()
	Stop()

	//从当前连接获取原始的socket TCPConn GetTCPConnection() *net.TCPConn //获取当前连接ID
	GetConnId() uint32
	GetTCPConnection() *net.TCPConn
	RemoteAddr() net.Addr
	// 直接将msg数据发给远程的tcp客户端
	SendMsg(msgId uint32, data []byte) error

}

// 定义一个统一处理连接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error  // HandFunc 是定义的新类型，类型是函数