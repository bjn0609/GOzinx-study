package znet

type Message struct {
	Id      uint32 //消息Id
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

//创建message消息包
func NewMsgPackage(id uint32,data []byte) *Message{
	return &Message{
		Id: id,
		DataLen: uint32(len(data)),
		Data: data,
	}
}


//获取消息ID
func(m *Message) GetMsgId() uint32{
	return m.Id
}
//获取信息长度
func(m *Message) GetMsgLen() uint32{
	return m.DataLen
}
//获取信息内容
func(m *Message) GetData() []byte{
	return m.Data
}

//设置消息ID
func(m *Message) SetMsgId(id uint32){
	m.Id = id
}
//设置消息内容
func(m *Message) SetData(data []byte){
	m.Data = data
}
//设置消息长度
func(m *Message) SetDataLen(len uint32){
	m.DataLen = len
}