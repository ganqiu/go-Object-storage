// 提供接口服务
package Api

import (
	"data/Api/heartbeat"
	"data/Api/locate"
	"data/Api/objects"
	"log"
	"net/http"
	"os"
)

// 监听apiserver的信息并处理请求
func main() {
	go heartbeat.ListenHeartbeat() //开启监听函数
	http.HandleFunc("/objects/", objects.Handler)
	//分别处理以/locate/和/objects/为开头的函数
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
	//监听端口（由系统环境变量LISTEN_ADDRESS定义）运行后将始终监听，非正常情况将错误返回并打印错误信息退出
}
