// 向apiServers发送心跳消息，并创建临时消息队列发送一消息（正文为需要定位的对象，返回临时队列的名字）
package heartbeat

import (
	"data/rabbitmq"
	"os"
	"time"
)

// 开始监听
func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER")) //创建结构体
	defer q.Close()
	for { //无限循环
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS")) //调用Publish方法发送监听地址
		time.Sleep(5 * time.Second)                          //一定时间后再次发送心跳消息
	}
}
