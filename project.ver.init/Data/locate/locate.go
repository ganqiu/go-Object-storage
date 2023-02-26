// 提供PUT和GET方法
package locate

import (
	"data/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)    //获取磁盘文件名
	return !os.IsNotExist(err) //判断文件名是否存在

} //定位文件是否存在，并返回True/False
func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ-server")) //创建结构体
	defer q.Close()
	q.Bind("dataServers") //绑定dataServer exchange
	c := q.Consume()      //返回一个channel
	for msg := range c {  //遍历channel来接收消息
		//消息内容：接口服务发送的定位对象名->经过json编码对象带有双引号
		object, e := strconv.Unquote(string(msg.Body)) //去除双引号并返回字符串结果
		if e != nil {
			panic(e) //内建函数panic停止当前Go程的正常执行。当函数调用panic时，该函数的正常执行就会立刻停止。
			// 函数中defer的所有函数先入后出执行后，函数返回给其调用者。调用者如同函数一样行动，层层返回，
			//直到该Go程中所有函数都按相反的顺序停止执行。之后，程序被终止，
			//而错误情况会被报告，包括引发该恐慌的实参值，此终止序列称为恐慌过程。
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/OBJECTS/" + object) /*调用Locate函数检查文件是否存在*/ {
			//Getenv函数作用：检索并返回名为key的环境变量的值。如果不存在该环境变量会返回空字符串。
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
			//调用Send向消息的发送方返回服务节点监听地址，表明对象在该服务节点
		}
	}
}
