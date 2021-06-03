package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

//模拟客户端
func main(){
	fmt.Println("client1 start...")

	time.Sleep(1 *time.Second)

	//1连接远程服务器，得到一个Conn连接
	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil{
		fmt.Println("client start err,exit!",err)
		return
	}

	for{
	//发送封包的message消息 msgID:0
	dp := znet.NewDataPack()
	binaryMsg,err := dp.Pack(znet.NewMsgPackage(1,[]byte("Zinx client1 Test Message")))	
	if err != nil {
		fmt.Println("Pack error:",err)
		return
	}	
	if _,err := conn.Write(binaryMsg);err != nil {
		fmt.Println("write error",err)
		return
	}

	//服务器应该回复一个message数据 msgID:1 pingpingping

	//1先读取流中的head部分 得到ID 和dataLen

	binaryHead := make([]byte,dp.GetHeadLen())
	if _,err := io.ReadFull(conn,binaryHead);err != nil{
		fmt.Println("read head error",err)
		break
	}
	//将二进制head拆包到msg结构体中
	msgHead,err := dp.Unpack(binaryHead)
	if err != nil {
		fmt.Println("client Unpack msgHead error",err)
		break
	}

	if msgHead.GetMsgLen() > 0 {
		//2再根据DataLen进行第二次读取 将data读出来
		msg := msgHead.(*znet.Message)
		msg.Data = make([]byte, msg.GetMsgLen())

		if _,err := io.ReadFull(conn,msg.Data);err != nil{
			fmt.Println("read msg data error,",err)
			return
		}
		fmt.Println("———>Recv Server Msg : ID = ",msg.Id,",len = ",msg.DataLen,"data = ",string(msg.Data))
	}

	

	//cpu阻塞
	time.Sleep(1*time.Second)
	}
}