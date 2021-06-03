package znet

import (
	"zinx/ziface"
)


type Reqeust struct {
	//已经和客户端建立号的连接Reqeust
	conn ziface.IConnection

	//客户端请求的数据
	msg ziface.IMessage
}

//得到当前连接
func(r *Reqeust) GetConnection() ziface.IConnection{
	return r.conn
}

//得到请求的消息数据
func(r *Reqeust) GetData() []byte{
	return r.msg.GetData()
}

func(r *Reqeust) GetMsgID() uint32{
	return r.msg.GetMsgId()
}