package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	//服务器名称Server
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息关系模块，用来绑定MSGID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	//该server的连接管理器
	ConnMgr ziface.IConnManager
	//该server创建链接之后自动调用HOOK函数
	OnConnStart func(conn ziface.IConnection)
	//该server销毁链接之前自动调用HOOK函数
	OnConnStop func(conn ziface.IConnection)
}



//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx]Server Name : %s, listenner at IP: %s, Port:%d is starting\n",
		utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx]Version : %s, MaxConn : %d, MaxPacketSize : %d\n",
		utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSize)
	
	go func() {
		//0 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()

		//1获取TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve TCP Addr error:", err)
			return
		}
		//2监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")
		var cid uint32 = 0
		//3阻塞地等待客户端连接，处理客户端连接服务（读写）
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断，超过最大连接数量，关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn{
				//TODO 给客户端相应超出最大连接的错误包
				fmt.Println("====>Too Many Connection MaxConn = ",utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn绑定，得到连接模块
			dealConn := NewConnection(s,conn,cid,s.MsgHandler)
			cid ++

			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	//TODO 将服务器的资源，状态或以及开辟的连接信息进行停止或回收
	fmt.Println("[STOP] Zinx server name ",s.Name)
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Server() {
	//启动Server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}
func (s *Server)AddRouter(msgID uint32,router ziface.IRouter){
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager{
	return s.ConnMgr
}

//初始化server的方法
func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:NewConnManager(),
	}
	return s
}


//注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart (hookFunc func(connection ziface.IConnection)){
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop (hookFunc func(connection ziface.IConnection)){
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart (conn ziface.IConnection){
	if s.OnConnStart != nil{
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop (conn ziface.IConnection){
	if s.OnConnStop != nil{
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStop(conn)
	}
}
