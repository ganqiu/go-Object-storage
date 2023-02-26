// 作为数据节点存在，提供数据的存储功能
package data

import (
	"data/Api/objects"
	"data/Data/heartbeat"
	"data/Data/locate"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objetcs/", objects.Handler)
	//访问本机HTTP且URL以/objetcs/开头时将由objects.Handler处理，除此之外一律返回404
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
	//监听端口（由系统环境变量LISTEN_ADDRESS定义）运行后将始终监听，非正常情况将错误返回并打印错误信息退出
}
