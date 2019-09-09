package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/prometheus/common/log"
	"net"
	"net/http"
	"net/http/pprof"
	"test/test1/conf"
	http2 "test/test1/httpclient"
	"test/test1/httpser"
	"test/test1/model"
	"test/test1/utils"
	"time"
)

var P = flag.Int("p", 8888, "p")
var Load = flag.String("load", "", "8881,8882,8883")
var Ps = ""

func Run(p *int){
	fmt.Println("服务端创建成功:",*p)
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "", *p),
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



	mux.Group(func(r chi.Router) {
		r.Use(Forward)
		r.Use(Distribute)
		r.Post("/http/client/register", Client{}.Register)
		r.Post("/http/client/delete", Client{}.Delete)
		r.Post("/http/client/send", Client{}.Send) //已经遗弃
		r.Get("/http/client/get", Client{}.Get)
		r.Post("/http/client/gomain", Colony{}.Gomain)
		r.Post("/http/client/in/user", User{}.InUser)	//用户ip地址
		r.Get("/http/client/user", User{}.SelectUser)
		r.Post("/http/client/de/user", User{}.DeleteUser)
	})

	return mux
}

func init() {
	client = Client{Addrs: make(map[string]string,0)}
	htpp = http2.Http{Timeout: 10}
	flag.Parse()
}

var htpp http2.Http

var client Client

var count int

type Client struct {
	Addrs    map[string]string
	HttpRule int64 // 1.是平均访问  2.是均衡访问 3.
}
type Colony struct {

}

//转发
func Forward(next http.Handler)http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//如果等于true的话则进行转发
		if request.Header.Get("forward") != "true"{
			add(writer,request)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

//转发，转发请求进行负载
func add(writer http.ResponseWriter, request *http.Request){
	//删除其中一个地址和重新来过
	addNil:
	//每次进行加1 这样就会不停的轮询
	count++

	if len(client.Addrs) == 0{
		utils.WriterJson(writer, map[string]interface{}{"err": "The service is not registered", "code": 500, "data": nil})
		return
	}
	count := count%len(client.Addrs)
	//没有问题  用余数求出已经到哪个服务
	add := ""
	var i = 0
	for _,v := range client.Addrs{
		if i == count{
			add = v
		}
		i++
	}
	//input.Route

	//ping不通就删除
	if _, err := net.DialTimeout("tcp", add, 3*time.Second); err != nil {
		delete(client.Addrs,add)
		if len(client.Addrs) == 0{
			utils.WriterJson(writer, map[string]interface{}{"err": "client.Addrs is nil", "code": 500, "data": nil})
			return
		}
		goto addNil
	}

	httpser.Forward(writer,request,add)
}

//目前做个ip地址的校验
func Distribute(next http.Handler)http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//验证ip地址
		if !ipCheck(writer,request){
			return
		}

		var par string
		//拼接参数
		for k, v := range request.Form {
			par = fmt.Sprintf("%v%v",par ,fmt.Sprintf("%v=%v&",k,v))
		}
		par = fmt.Sprintf("%v%v",par,"wjwp=wjwp")
		//token
		token := request.Header.Get("x-auth-token")
		//fmt.Println("token:",token)
		request = request.WithContext(context.WithValue(request.Context(), "apiKey", token))
		check := request.Header.Get("check")
		//如果已经验证一遍则跳过
		if check == "false1"{
			return
		}
		//发送给别的服务同步数据
		/*for _,v := range strings.Split(*Load,","){
			vInt,_ := strconv.Atoi(v)
			if vInt != *P{

				fmt.Println("Load:",fmt.Sprintf("http://127.0.0.1:%v",vInt))
				//是否ping通了
				//请求必须为：http://127.0.0.1:9090   会截取后面的
				if _, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%v",vInt), 3*time.Second); err != nil {
					fmt.Println(err)
					continue
				}

				//目的地都发送一遍
				_, err := http2.HttpClient.PostClient(fmt.Sprintf("http://127.0.0.1:%v%v",vInt,request.RequestURI),par, false, map[string]string{"x-auth-token":token,"check":"false1"})
				if err != nil {
					fmt.Println("needp:",vInt,err)
					//panic(nil)
				}
			}
		}
		//访问对面的程序重启服务 访问接口
		/*_, err := htpp.Post(fmt.Sprintf("http://127.0.0.1:%v/http/client/gomain",needP),fmt.Sprintf("addrs=%v",Ps) , false)
		if err != nil {
			fmt.Println("needp:",needP,err)
			panic(nil)
		}*/
		next.ServeHTTP(writer, request)
	})
}


func getContext(writer http.ResponseWriter, request *http.Request) *Contextr {
	return &Contextr{
		Token: request.Context().Value("apiKey").(string),
		Request:  request,
		Writer:   writer,
	}
}

var Contextrp Contextr

type Contextr struct {
	Token string
	Request *http.Request
	Writer  http.ResponseWriter
}


//测试通过
func ipCheck(writer http.ResponseWriter, request *http.Request) bool{
	//fmt.Println("RealIP:",httpser.RealIP(writer,request))
	//fmt.Println("Platform:",httpser.GetPlatform(request))
	//开启了ip验证
	if conf.LoadConf.PrescribeIp == true{
		//数据库查询ip
		users , err := model.Db.UserSelectToParam(map[string]string{"ip":httpser.RealIP(writer,request)})
		if err != nil{
			log.Error("ipCheck.UserSelectToParam")
			utils.WriterJson(writer, map[string]interface{}{"err": err.Error(), "code": 400, "data": nil})
			return false
		}
		if len(users) < 1 {
			utils.WriterJson(writer, map[string]interface{}{"err": "ip地址没有注册", "code": 400, "data": nil})
			return false
		}
		//已经使用过的数量大于允许最大数量就退出  不能等于-1 -1就是永久访问
		if users[0].UseIpVisitCount >= users[0].IpVisitCount && users[0].IpVisitCount != -1{
			utils.WriterJson(writer, map[string]interface{}{"err": "ip当前访问数量已用完", "code": 400, "data": nil})
			return false
		}
		sql := "update user set use_ip_visit_count = use_ip_visit_count + 1 where ip = ?"
		err = model.Db.Update(sql,httpser.RealIP(writer,request))
		if err!=nil{
			log.Error("ipCheck.Update")
			utils.WriterJson(writer, map[string]interface{}{"err": err.Error(), "code": 400, "data": nil})
			return false
		}
	}
	return true
}