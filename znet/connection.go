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

	handleAPI    ziface.HandFunc
	ExitBuffChan chan bool // 告知该连接已经退出/停止的channel
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callback_api,
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
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("revc buff err:", err)
			c.ExitBuffChan <- true
			continue
		}
		if err := c.handleAPI(c.Conn, buf, cnt); err!=nil{
			fmt.Println("connId", c.ConnID, "handle is error:", err)
			c.ExitBuffChan <- true
			return
		}
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

	c.ExitBuffChan <- true  // 通知缓冲队列读数据的业务，该连接已经关闭

	close(c.ExitBuffChan)  // 关闭一个通道会触发所有还在等待该通道数据的 case 语句执行
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