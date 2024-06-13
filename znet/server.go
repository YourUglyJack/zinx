package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/ziface"
)

type Server struct {
	Name      string // server name
	IPVersion string // tcp4 or other
	IP        string
	Port      int
}

// ================== 回显业务 ====================
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buff err:", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// ================== 实现ziface.IServer 里的全部接口方法 ====================

func (s *Server) Start() {
	fmt.Printf("[Start] Server listenner at IP: %s, Port: %d, is starting...\n", s.IP, s.Port)

	go func() {

		// 获取TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}

		// 监听
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}

		fmt.Println("start Zinx server ", s.Name, "succ, now listening...")

		// TODO 自动生成id的方法
		cid := uint32(0)

		// 启动server网络连接业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// TODO 设置服务器最大连接

			// TODO 处理新连接 请求业务的方法， handle和conn是绑定的
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid ++

			// 启动当前连接的业务
			go dealConn.Start()

			// // 暂时做一个最大512字节的echo的服务
			// go func() {
			// 	for {
			// 		buf := make([]byte, 512)
			// 		cnt, err := conn.Read(buf)
			// 		if err != nil {
			// 			fmt.Println("recv buf err:", err)
			// 			continue
			// 		}
			// 		// echo, func (c *TCPConn) Write(b []byte) (int, error)
			// 		if _, err := conn.Write(buf[:cnt]); err != nil {
			// 			fmt.Println("Write back buf err:", err)
			// 			continue
			// 		}
			// 	}
			// }()
		}

	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server ", s.Name)

	// TODO Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

func (s *Server) Server() {
	s.Start()

	// 阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      7778,
	}

	return s
}
