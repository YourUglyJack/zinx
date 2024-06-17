package main

import (
	"fmt"
	"io"
	"net"
	"zinx/znet"
)

func main() {
	listenner, err := net.Listen("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			break
		}

		go func (conn net.Conn)  {
			dp := znet.NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head err:", err)
					break
				}

				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack err:", err)
					return
				}

				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*znet.Message)  // imessage -> message
					msg.Data = make([]byte, msg.GetDataLen())

					// 根据datalen从io中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("unpack data err:", err)
						return
					}
					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)

	}
}