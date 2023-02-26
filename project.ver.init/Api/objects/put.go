package objects

import (
	"data/Api/heartbeat"
	"data/Api/objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// 获取<object_name>给object，将r.Body和Object作为参数调用storeObject，并返回一个HTTP错误代码
// 将错误代码写入HTTP响应
// storeObject也会返回一个error，当error不等于nil就打印在log中
func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}

// 以object为参数调用putStream生成一个指向Putstream结构体的指针，该结构体实现了Write方法，通过使用
// io.Copy将请求正文写入stream，再调用stream.Close（）关闭stream流，Close()会返回一个error来通知错误
// 如果有错误会返回错误代码500
func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

// 调用函数获取一个随机数据服务节点，如果空字符创则返回一个空指针和"cannot find any dayaServer"的错误
// store会返回StatusServiceUnavailable，客户端会收到503error，不为空则生成objectstream.PutStream指针并返回
func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dayaServer")
	}
	return objectstream.NewPutStream(server, object), nil
}
