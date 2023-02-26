// 该包不会访问本地磁盘而是转发给数据服务进行本地访问
package objects

import (
	"data/Api/locate"
	"data/Api/objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// 获取<object.name>，然后调用getStream生成一个io.Reader的stream，出现错误打印log并返回404，否则将stream写入HTTP响应的正文
func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

// object代表对象的名字，先调用Locate定位，如果返回空字符串说明定位失败
// 否则调用objectstream.NewGetStream并返回结果
func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("objects %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
