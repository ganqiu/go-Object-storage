// 将http函数的调用转换为读写流形式，便于处理
package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

// url获取数据流HTTP服务地址，然后调用http.GET获取HTTP响应，并返回r(类型为*http.Response),其中StatusCode（HTTP响应的错误代码）
// Body则被用于读取HTTP响应正文的io.Reader
// 从io.Reader获取响应正文
// 将r.Body作为新的reader返回Getstream
func newGetStream(url string) (*GetStream, error) {
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return &GetStream{r.Body}, nil
}

// server和object拼成一个url传给newGetStream，隐藏url细节使使用者只需提供服务节点地址和对象名就可以读取对象
func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invaild server %s objects %s", server, object)
	}
	return newGetStream("http://" + server + "/objects/" + object)
}

// 读取reader成员，使Getstream实现io.Reader接口
func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
