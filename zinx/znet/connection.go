package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前conn隶属于哪个server
	TcpServer ziface.Iserver

	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前连接的状态
	isClose bool
  
	//告知当前链接已经退出的channel   Reader告诉Writer退出
	ExitChan chan bool

	//无缓冲的通道 用于读写goroutine之间的消息通讯
	msgChan chan []byte

	//消息的管理MSGID 和对应的处理业务api关系
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{} 

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

//初始化连接模块的方法
func NewConnection(server ziface.Iserver,conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandler: msgHandler,
		isClose:   false,
		msgChan: make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property: make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//连接的读数据的业务
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine us running...]")
	defer fmt.Println("[Reader is exit],connID = ", c.ConnID, "remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

		//创建拆包解包对象
		dp := NewDataPack()

		//读取客户端的MSG HEAD二进制流 8个字节
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(c.GetTCPConnection(),headData);err != nil{
			fmt.Println("read msg head error ",err)
			break
		}

		//拆包 得到msgID和msgDatalen放在msg消息在
		msg,err := dp.Unpack(headData)
		if err != nil{
			fmt.Println("unpack error ",err)
			break
		}

		//根据datalen 再次读取data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0{
			data = make([]byte, msg.GetMsgLen())
			if _,err := io.ReadFull(c.GetTCPConnection(),data);err != nil{
				fmt.Println("read msg data error ",err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的request请求数据
		req := Reqeust{
			conn: c,
			msg: msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启工作池，将消息发送给worker工作池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}else{
			//从路由中，找到注册绑定的Conn对应的router调用
			//根据绑定好的MSGID 找到对应处理api业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//写消息的goroutine，专发送给客户端消息的方法
func (c *Connection) StartWriter(){
	fmt.Println("[Write Goroutine is running...]")
	defer fmt.Println("[conn Writer exit!]",c.RemoteAddr().String())

	//不断阻塞的等待chan的消息，进行写给客户端
	for {
		select{
		case data := <-c.msgChan:
			//有数据写给客户端
			if _,err := c.Conn.Write(data);err != nil{
				fmt.Println("Send data error,",err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出，writer也要退出
			return
		}

	}
}


//启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start ... ConnID = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接的写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应HOOK函数
	c.TcpServer.CallOnConnStart(c)
}

//停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()..ConnID = ", c.ConnID)

	//若连接已经关闭Stop
	if c.isClose {
		return
	}
	c.isClose = true

	//按照开发者传递进来的 销毁链接之前需要调用的处理业务，执行对应HOOK函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从connmgr中摘除
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供sendmsg方法，将发给客户端的数据 先封包 再发送
func (c *Connection) SendMsg(msgId uint32,data []byte)error{
	if c.isClose{
		return errors.New("Connection cloesd when send msg")
	}

	//将data封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	binaryMsg,err := dp.Pack(NewMsgPackage(msgId,data))
	if err != nil{
		fmt.Println("Pack error msg id = ",msgId)
		return errors.New("Pack error msg")
	}

	//数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

//设置连接属性
func (c *Connection) SetProperty(key string,value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个连接属性
	c.property[key] = value
}



//获取连接属性
func (c *Connection) GetProperty(key string) (interface{},error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value,ok := c.property[key]; ok {
		return value,nil
	}else{
		return nil,errors.New("no property found")
	}
}

//移除连接属性
func (c *Connection) RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property,key)
}


