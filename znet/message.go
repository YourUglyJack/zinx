package znet

type Message struct { // 为了解决TCP粘包问题
	Id      uint32
	DataLen uint32
	Data    []byte
}

// 创建一个msg消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

func (msg *Message) SetMsgId(id uint32) {
	msg.Id = id
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}


