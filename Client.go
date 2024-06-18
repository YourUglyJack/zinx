package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
	模拟客户端
*/
func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn,err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.5 Client Test")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("Write err:", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head err:", err)
			break
		}
		
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack err:", err)
			}
			fmt.Println("==> Recv MsgID:", msg.Id, ", len:", msg.DataLen, ", data:", string(msg.Data))
		}
		// _, err := conn.Write([]byte("Zinx V0.3"))
		// if err !=nil {
		// 	fmt.Println("write error err ", err)
		// 	return
		// }

		// buf :=make([]byte, 512)
		// cnt, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("read buf error ")
		// 	return
		// }

		// fmt.Printf(" server call back : %s, cnt = %d\n", buf,  cnt)

		time.Sleep(1*time.Second)
	}
}