package main

import (
	"fmt"
	"net"
	"zinx/znet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("client conn err:", err)
		return
	}

	dp := znet.NewDataPack()

	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err:", err)
		return
	}

	msg2 := &znet.Message{
		Id:      1,
		DataLen: 6,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!'},
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err:", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	select {}

}
