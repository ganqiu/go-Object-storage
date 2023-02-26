package objects

import "net/http"

// 写入HTTP响应
// ResponseWriter写入HTTP响应，WriteHeader响应错误代码，Write写HTTP响应正文
// r类型：*http.Request->代表当前处理的HTTP请求
// Handler检查HTTP请求方法，并调用对应函数，如果都不成立就返回代码405
func Handler(w http.ResponseWriter, r *http.Request) {

	m := r.Method //Method记录请求方法
	if m == http.MethodPut {
		put(w, r) //存储到本地硬盘
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
