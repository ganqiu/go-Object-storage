package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

// wirter实现Write方法，c实现将传输数据发生的错误传回主线程
type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

// 生成一个PutStream结构体，运用了io.Pipe创建了一对reader和writer，类型为*io.PipeReader和*io.PipeWriter
// 二者管道互联，写入writer可以从reader读出
// 因此可以以数据流形式实现Put请求，并通过类型为http.Client的变量中读取PUT的内容，
// 从而实现同时写入writer实现Putstream的Write方法
// 又因为管道读写堵塞，故在协程中调用client.Do的方法，该方法返回错误代码和error，如果error不为空则将错误
// 发送到协程上中，如果error为空但是HTTP错误代码不为200也认为是一种错误，并在之后被PUTSTREAM.Close读取
func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		client := http.Client{}
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataServer erturn http code%d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

// 用于写入writer实现io.write接口
func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// 用于关闭writer，让管道另一端reader读取到io.EOF,否则client.Do将阻塞无法返回
func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}
