package znet

import (
	"fmt"
	"zinx/utils"
	"strconv"
	"zinx/ziface"
)

//消息处理模块的实现
type MsgHandle struct{
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32] ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

//初始化 创建MSGHANDLE方法
func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		Apis:make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,//从全局配置获取
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest){
	//1从request中找到msgID
	handler,ok := mh.Apis[request.GetMsgID()]
	if !ok{
		fmt.Println("api msgID = ",request.GetMsgID()," is NOT FOUND! Need Register!")
	}

	//2根据msgID 调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32,router ziface.IRouter){
	//1判断当前msg绑定的API处理方法是否存在
	if _,ok := mh.Apis[msgID];ok{
		//id已经注册
		panic("repeat api,msgID = " + strconv.Itoa(int(msgID)))
	}

	//2添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ",msgID,"succ!")
}
 
//启动一个worker工作池  开启工作池的动作只能发生一次  只有一个工作池
func (mh *MsgHandle) StartWorkerPool(){
	//根据WorkerPoolSize 分别开启worker 每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 给当前worker对应的channel消息队列 开辟空间 第0个worker用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest,utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i,mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int,taskQueue chan ziface.IRequest){
	fmt.Println("Worker ID = ",workerID," is started...")

	//不断阻塞等待对应消息队列的消息
	for {
		select{
			//若有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务
			case request := <- taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest){
	//1 将消息平均分配给不同的worker
	//根据客户端建立的ConnID分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ",request.GetConnection().GetConnID(),
	"request MsgID = ",request.GetMsgID()," to workerID = ",workerID)

	//2 将消息发送给对应worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}