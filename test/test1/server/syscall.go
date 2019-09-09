package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Control(p int,load string) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	fmt.Println("receive signal `%s`\n", <-sig)
	//需要访问的端口
	needP := 0

	//循环负载端口，拿出端口进行重新加载服务
	for _,v:= range strings.Split(load,","){
		pOne , err := strconv.Atoi(v)
		if err != nil{
			fmt.Println(err)
			panic(nil)
		}
		if pOne != p{
			//是否ping通了
			//请求必须为：http://127.0.0.1:9090   会截取后面的
			if _, err := net.DialTimeout("tcp", fmt.Sprintf(":%v",pOne), 3*time.Second); err == nil {
				needP = pOne
			}
		}
	}
	if needP == 0{
		fmt.Println("程序终止")
		panic(nil)
	}
	//访问对面的程序重启服务 访问接口
	_, err := htpp.Post(fmt.Sprintf("http://:%v/http/client/gomain",needP),fmt.Sprintf("addrs=%v",Ps) , false,nil)
	if err != nil {
		fmt.Println("needp:",needP,err)
		panic(nil)
	}
	fmt.Println(fmt.Sprintf("http://:%v/http/client/gomain",needP))
	panic(nil)
}