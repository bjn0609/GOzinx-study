package znet
import(
	"zinx/ziface"
)

//实现router时，先嵌入这个baserouter基类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct{}

//都为空，是因为有的router不希望有PreHandle PostHandle这两个业务
//所以Router全部继承BaseRouter的好处就是 不需要实现PreHandle PostHandle

//在处理Conn业务之前的钩子方法
func (br *BaseRouter) PreHandle(request ziface.IRequest){}

//在处理Conn业务的主钩子方法
func (br *BaseRouter) Handle(request ziface.IRequest){}

//在处理Conn业务之后的钩子方法
func (br *BaseRouter) PostHandle(request ziface.IRequest){}
