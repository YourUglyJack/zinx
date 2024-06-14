package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	Router ziface.IRouter

	// handleAPI    ziface.HandFunc
	ExitBuffChan chan bool // 告知该连接已经退出/停止的channel
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit...")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("revc buff err:", err)
			c.ExitBuffChan <- true
			continue
		}

		// 在conn读取完客户端数据之后，将数据和conn封装到一个Request中，作为Router的输入数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 然后开启一个goroutine去调用给Zinx框架注册好的路由业务
		go func(request ziface.IRequest) {
			fmt.Println(c.Router)
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		// if err := c.handleAPI(c.Conn, buf, cnt); err!=nil{
		// 	fmt.Println("connId", c.ConnID, "handle is error:", err)
		// 	c.ExitBuffChan <- true
		// 	return
		// }
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO

	c.Conn.Close() // 关闭socket链接

	c.ExitBuffChan <- true // 通知缓冲队列读数据的业务，该连接已经关闭

	close(c.ExitBuffChan) // 关闭一个通道会触发所有还在等待该通道数据的 case 语句执行
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
