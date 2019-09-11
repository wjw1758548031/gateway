package httpclient

import (
	"github.com/prometheus/common/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handle struct {
	host string
	port string
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse("http://" + this.host + ":" + this.port)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func StartServer(port, localPort string) {
	//被代理的服务器host和port 本地端口 必须要填写一个ip地址
	h := &handle{host: "127.0.0.1", port: localPort}
	//监控的端口
	err := http.ListenAndServe("127.0.0.1:"+port, h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
