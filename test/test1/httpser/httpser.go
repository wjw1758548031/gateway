package httpser

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)


func  RealIP(writer http.ResponseWriter, request *http.Request) string {
	if s := request.Header.Get(HeaderXForwardedFor); s != "" {
		return strings.Split(s, ", ")[0]
	}
	if s := request.Header.Get(HeaderXRealIP); s != "" {
		return s
	}
	ra, _, _ := net.SplitHostPort(request.RemoteAddr)
	return ra
}

type Platform int

const (
	IOS = Platform(iota)
	Android
	Web
)


func GetPlatform(r *http.Request) Platform {
	client := r.Header.Get("client")
	if strings.Contains(client, "ios") {
		return IOS
	}
	if strings.Contains(client, "android") {
		return Android
	}
	return Web
}

func (p Platform) String() string {
	switch p {
	case IOS:
		return "IOS"
	case Android:
		return "ANDROID"
	default:
		return "WEB"
	}
}


func Forward(w http.ResponseWriter, req *http.Request,host string) {
	targetHost := &TargetHost{
		Host:    host,//127.0.0.1:8888
	}
	HostReverseProxy(w, req, targetHost)
}
func HostReverseProxy(w http.ResponseWriter, req *http.Request, targetHost *TargetHost) {
	host := ""
	if targetHost.IsHttps {
		host = host + "https://"
	} else {
		host = host + "http://"
	}
	remote, err := url.Parse(host + targetHost.Host)
	if err != nil {
		fmt.Println("err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
}
type TargetHost struct {
	Host    string
	IsHttps bool
	CAPath  string
}
