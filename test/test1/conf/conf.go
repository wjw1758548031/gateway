package conf

import (
	"gopkg.in/ini.v1"
	"log"
)
var filepath ="load.conf"

type Config struct {   //配置文件要通过tag来指定配置文件中的名称
	PrescribeIp bool  `ini:"prescribe_ip"` //是否允许开启指定ip false 或 true
}

var LoadConf Config


func init() {
	loadConf, err := ReadConfig(filepath) //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Fatal(err)
		return
	}
	LoadConf = loadConf
	log.Println(LoadConf)
}

//读取配置文件并转成结构体
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)   //加载配置文件
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)   //解析成结构体
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
