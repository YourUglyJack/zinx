package ziface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte  // data of request
	GetMsgId() uint32
}