package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	MsgHandle ziface.IMsgHandle
	// 告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	// 无缓冲通道，用于读写两个goroutine之间的消息通信
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandle:    msgHandle,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
	return c
}

// 处理conn读数据的goroutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit...")
	defer c.Stop()

	for {
		// 创建拆包解包对象
		dp := NewDataPack()

		// 读取客户端的msghead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:", err)
			c.ExitBuffChan <- true
			continue
		}

		// 拆包，得到msg，msg里有msgid以及datalen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			c.ExitBuffChan <- true
			continue
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg err", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 在conn读取完客户端数据之后，将数据和conn封装到一个Request中，作为Router的输入数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer eixt]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err:", err)
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	go c.StartWriter()

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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack err, msg id:", msgId)
		return errors.New("pack err")
	}

	c.msgChan <- msg

	return nil
}
