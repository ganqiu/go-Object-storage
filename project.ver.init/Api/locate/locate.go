// 用于定位的locate接口来验证架构，当接受客户端请求时会像数据服务发送定位消息并返回节点/404
package locate

import (
	"data/rabbitmq"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 处理HTTP请求，不为GET方法则返回405，方法为GETz则以<object_name>作为参数来定位。Locate
// 返回为空则说明定位失败，返回404
// 不为空则发挥地址作为HTTP响应正文
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

// 接受一个string的参数name（定位对象的名字），并创建一个新的消息队列向dataServers exchange群发对象
// 名字定位消息
// 协程启动匿名函数并设置1s后关闭->设置超时机制
// 1s后无反馈则消息队列关闭接收消息长度为0，返回空字符串，有消息则返回节点监听地址
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// 检查Locate是否为空字符串来判定对象是否存在
func Exist(name string) bool {
	return Locate(name) != ""
}
