package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("recv from client:msgID = ", request.GetMsgID(),
		",data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//hello zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping..ping..ping
	fmt.Println("recv from client:msgID = ", request.GetMsgID(),
		",data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello welcome to zinx!!"))
	if err != nil {
		fmt.Println(err)
	}
}

//创建链接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection){
	fmt.Println("===> DoConnectionBegin is Called...")
	if err := conn.SendMsg(202,[]byte("DoConnection BEGIN")); err != nil{
		fmt.Println(err)
	}

	//给当前的连接设置一些属性
	fmt.Println("Set conn property...")
	conn.SetProperty("Name","——Aceld——")
	conn.SetProperty("Home","https://github.com/aceld")
}

//链接断开之前执行的钩子函数
func DoConnectionLost(conn ziface.IConnection){
	fmt.Println("===> DoConnectionLost is Called...")
	fmt.Println("conn ID = ",conn.GetConnID()," is Lost...")

	//获取链接属性
	if name,err := conn.GetProperty("Name");err == nil{
		fmt.Println("Name = ",name)
	}

	if home,err := conn.GetProperty("Home");err == nil{
		fmt.Println("Home = ",home)
	}
}

func main() {
	//1创建一个server句柄，使用Zinx的API
	s := znet.NewServer("[zinx V1.0]")

	//2 注册链接HOOK钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3给当前框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//4启动server
	s.Server()
}
