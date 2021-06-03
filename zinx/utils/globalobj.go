package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切有关框架的全局参数，供其他模块使用
//一些参数是可以通过zinx.json由用户配置

type GlobalObj struct{
	//Server
	TcpServer ziface.Iserver	//当前Zinx全局的server对象
	Host string					//当前服务器主机监听的IP
	TcpPort int					//当前服务器主机监听的端口
	Name string					//当前服务器名称

	//Zinx
	Version string 				//框架版本号
	MaxConn int					//当前服务器主机允许的最大连接数
	MaxPackageSize uint32		//当前框架数据包的最大值
	WorkerPoolSize uint32		//当前业务工作worker池的协程数量
	MaxWorkerTaskLen uint32		//框架允许用户最多开辟多少个Worker
}

//定义全局的对外Globalobj
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload(){
	data,err := ioutil.ReadFile("conf/zinx.json")
	if err != nil{
		panic(err)
	}
	//json文件数据解析到struct中
	err = json.Unmarshal(data,&GlobalObject)
	if err != nil{
		panic(err)
	}
}


//初始化当前的GlobalObject
func init(){
	//若配置文件没加载 默认值
	GlobalObject = &GlobalObj{
		Name: "ZinxServerApp",
		Version: "V1.0",
		TcpPort: 8999,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,		//worker工作池的队列的个数
		MaxWorkerTaskLen: 1024,	//每个worker对应的消息队列的任务的数量最大值
	}

	//从conf/zinx.json加载用户自定义参数
	GlobalObject.Reload()
}