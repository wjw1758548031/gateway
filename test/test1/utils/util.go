package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/common/log"
	"net/http"
)

//获取参数
func Binding(r *http.Request, v interface{}) error {
	b := binding.Default(r.Method, r.Header.Get("Content-Type"))
	return b.Bind(r, v)
}

//输出json格式
func WriterJson(writer http.ResponseWriter, message map[string]interface{}) {
	e := json.NewEncoder(writer)
	err := e.Encode(message)
	if err != nil {
		log.Error(err)
	}
}

