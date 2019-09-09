package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	http2 "test/test1/httpclient"
	"test/test1/utils"
	"time"
)

//程序挂了重启服务
func (Colony) Gomain(writer http.ResponseWriter, request *http.Request) {
	input := &struct {
		Addrs string `form:"addrs" query:"addrs"`
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		fmt.Println("parameter error ")
		return
	}
	if input.Addrs == ""{
		fmt.Println("The address is empty")
		return
	}
	for _,v := range strings.Split(input.Addrs,","){
		//p ,_ :=strconv.Atoi(v)
		go http2.StartServer(v,strconv.Itoa(*P))
		//go Run(&p)
	}
	Ps = fmt.Sprintf("%v,%v",input.Addrs,Ps)
}


//查看服务内的地址
func (Client) Get(writer http.ResponseWriter, request *http.Request) {
	utils.WriterJson(writer, map[string]interface{}{"err": nil, "code": 200, "data": client.Addrs})
}


//发送对应的接口里在进行接受 //参数格式  id=1&name=zz 也可以直接序列化
func (Client) Send(writer http.ResponseWriter, request *http.Request) {
	//ctx := getContext(writer, request)
	input := &struct {
		Server string `form:"server" query:"server"`
		Route  string `form:"route" query:"route"`
		Mode   string `form:"mode" query:"mode"`
		Param  string `form:"param" query:"param"`
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		utils.WriterJson(writer, map[string]interface{}{"err": "parameter error ", "code": 400, "data": nil})
		return
	}
	if input.Route == "" || input.Mode == "" || input.Server == "" {
		utils.WriterJson(writer, map[string]interface{}{"err": "Parameters cannot be empty ", "code": 500, "data": nil})
		return
	}
	//删除其中一个地址和重新来过
	//addNil:
	//每次进行加1 这样就会不停的轮询
	count++

	if len(client.Addrs[input.Server]) == 0{
		utils.WriterJson(writer, map[string]interface{}{"err": "The service is not registered", "code": 500, "data": nil})
		return
	}
	//count := count%len(client.Addrs[input.Server])
	//没有问题  用余数求出已经到哪个服务
	//add := client.Addrs[input.Server][count]
	//input.Route

	//为空就删除
	/*if add == ""{
		//只剩最后一个时就有另外的规则了
		if len(client.Addrs[input.Server]) == 1{
			client.Addrs[input.Server] = []string{}
			return
		}
		client.Addrs[input.Server] = append(client.Addrs[input.Server][:count], client.Addrs[input.Server][count+1:]...)
		goto addNil
	}*/

	result := ""
	//看是什么请求
	/*if input.Mode == "post" {
		result, err = htpp.Post(fmt.Sprintf(add+input.Route),input.Param , false, map[string]string{"x-auth-token":ctx.Token})
		if err != nil {
			fmt.Printf("htpp.Post err=%v\n", err)
			utils.WriterJson(writer, map[string]interface{}{"err": "htpp.Post error", "code": 500, "data": nil})
			return
		}
	} else {
		utils.WriterJson(writer, map[string]interface{}{"err": "Please fill in the correct way", "code": 500, "data": nil})
		return
	}*/
	var data map[string]interface{}
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		fmt.Printf("unmarshal err=%v\n", err)
		utils.WriterJson(writer, map[string]interface{}{"err": "json.Unmarshal error ", "code": 500, "data": nil})
		return
	}
	utils.WriterJson(writer, map[string]interface{}{"err": nil, "code": 200, "data": data})
}

//删除注册地址
func (Client) Delete(writer http.ResponseWriter, request *http.Request) {
	input := &struct {
		Addr   string `form:"addr" query:"addr"`
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		utils.WriterJson(writer, map[string]interface{}{"err": "parameter error ", "code": 400, "data": nil})
		return
	}
	if input.Addr == ""  {
		utils.WriterJson(writer, map[string]interface{}{"err": "addr or server Can't be empty.", "code": 500, "data": nil})
		return
	}

	//查看当前服务里有没有对应的ip和端口  如果没有则直接进行退出
	if client.Addrs[input.Addr] == "" {
		utils.WriterJson(writer, map[string]interface{}{"err": "No service. ", "code": 500, "data": nil})
		return
	}
	//进行删除地址
	delete(client.Addrs, input.Addr)
	utils.WriterJson(writer, map[string]interface{}{"data": "success", "code": 200, "err": nil})
}

//Addr注册地址
func (Client) Register(writer http.ResponseWriter, request *http.Request) {
	input := &struct {
		Addr   string `form:"addr" query:"addr"`
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		utils.WriterJson(writer, map[string]interface{}{"err": "parameter error ", "code": 400, "data": nil})
		return
	}
	if input.Addr == ""  {
		utils.WriterJson(writer, map[string]interface{}{"err": "addr or server Can't be empty.", "code": 500, "data": nil})
		return
	}
	//判断下请求前是否加了http//
	if isOk, _ := regexp.MatchString("http://*", input.Addr); !isOk {
		utils.WriterJson(writer, map[string]interface{}{"err": "The address format error is, for example: http://:5000", "code": 500, "data": nil})
		return
	}

	//是否ping通了
	var conn net.Conn
	//请求必须为：http://127.0.0.1:9090   会截取后面的
	if conn, err = net.DialTimeout("tcp", strings.Split(input.Addr, "//")[1], 3*time.Second); err != nil {
		fmt.Println(err)
		utils.WriterJson(writer, map[string]interface{}{"err": input.Addr + ":No access ", "code": 500, "data": nil})
		return
	}
	conn.Close()
	client.Addrs[input.Addr] = input.Addr
	//把最多地址个数给到轮询位置 也相当于初始化了
	count = len(client.Addrs)
	utils.WriterJson(writer, map[string]interface{}{"data": "success", "code": 200, "err": nil})
}
