package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-chi/chi"
	"github.com/prometheus/common/log"
	"net/http"
	"net/http/pprof"
)

func main() {
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "", 9090),
		Handler: Router(chi.NewMux()),
	}
	s.ListenAndServe()
}

//路由
func Router(mux *chi.Mux) *chi.Mux {
	mux.Get("/debug/pprof/", pprof.Index)
	mux.Get("/debug/pprof/allocs", pprof.Index)
	mux.Get("/debug/pprof/block", pprof.Index)
	mux.Get("/debug/pprof/goroutine", pprof.Index)
	mux.Get("/debug/pprof/heap", pprof.Index)
	mux.Get("/debug/pprof/mutex", pprof.Index)
	mux.Get("/debug/pprof/threadcreate", pprof.Index)

	mux.Get("/debug/pprof/cmdline", pprof.Cmdline)
	mux.Get("/debug/pprof/profile", pprof.Profile)
	mux.Get("/debug/pprof/symbol", pprof.Symbol)
	mux.Get("/debug/pprof/trace", pprof.Trace)

	//mux.Post("/ieo/project/list", IeoProject{}.ProjectList)

	mux.Group(func(r chi.Router) {
		//	r.Use(AuthMember)
		// 项目列表接口（申购成功）
		r.Post("/http/client/register", Client{}.Register)
	})
	return mux
}

func init(){
	client = Client{}
}

var client Client

type Client struct {
	Addrs map[string][]string
	HttpRule int64  // 1.是平均访问  2.是均衡访问 3.
}

//Addr注册地址
func (Client) Register(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("1")
	e := json.NewEncoder(writer)
	err := e.Encode(map[string]interface{}{"swjw":"swjw"})
	if err != nil {
		log.Error(err)
	}
}

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
