// 监听数据服务节点发送的心跳消息
package heartbeat

import (
	"data/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time) //键值对，键为string，值为time.Time结构体
var mutex sync.Mutex

// 创建消息队列绑定apiSerrvers exchange，并通过 go channel监听心跳信息，将消息的正文内容(数据服务节点的监听地址作为键)，
// 将收到消息的时间存入dataServers
func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpireDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// 扫描dataServers,清除10s未接受到心跳信息的数据服务节点
func removeExpireDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// 遍历dataServers返回当前所有数据服务节点
// 队dataServers的读写需要锁的保护，防止多个协程并发读写导致错误，
// 优化：采用RWMutex读写锁
func GetDataServers() []string {
	mutex.Lock()
	//保护map并发读写 sync.mutex（互斥锁）
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// 随机挑出一个数据服务节点并返回，如果节点为空返回字符串
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
