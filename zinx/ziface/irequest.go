package ziface

//irequest接口  实际上把客户端请求的链接信息和请求数据包装到一个request中

type IRequest interface{
	//得到当前连接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	//得到请求的消息ID
	GetMsgID() uint32
}


