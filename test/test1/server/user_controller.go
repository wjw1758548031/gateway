package server

import (
	"fmt"
	"net/http"
	"test/test1/model"
	"test/test1/utils"
)

type User struct {

}

//删除用户
func (User) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	input := &struct {
		Id int64 `form:"id" query:"id"`
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		fmt.Println("parameter error ")
		return
	}
	if input.Id == 0{
		utils.WriterJson(writer, map[string]interface{}{"err": "parameter not null", "code": 500, "data": nil})
		return
	}
	sql := "delete from user where uid=?"
	err = model.Db.Delete(sql,input.Id)
	if err != nil{
		utils.WriterJson(writer, map[string]interface{}{"err": err.Error(), "code": 400, "data": nil})
		return
	}
	utils.WriterJson(writer, map[string]interface{}{"err": "", "code": 200, "data": "success"})
}

//新增用户
func (User) InUser(writer http.ResponseWriter, request *http.Request) {
	input := &struct {
		Name string `form:"name" query:"name"`	//名称
		Ip string `form:"ip" query:"ip"`		//ip地址
		IpVisitCount int64 `form:"ip_visit_count" query:"ip_visit_count"`  //ip允许访问的次数    -1永久访问
	}{}
	err := utils.Binding(request, input)
	if err != nil {
		fmt.Println("parameter error ")
		return
	}
	if input.Name == "" || input.IpVisitCount == 0 || input.Ip == ""{
		utils.WriterJson(writer, map[string]interface{}{"err": "parameter not null", "code": 500, "data": nil})
		return
	}
	sql := "INSERT INTO user(name, ip, ip_visit_count , use_ip_visit_count , created) values(?,?,?,?,datetime('now'))"
	err = model.Db.Insert(sql,input.Name,input.Ip,input.IpVisitCount, 0)
	if err != nil{
		utils.WriterJson(writer, map[string]interface{}{"err": err.Error(), "code": 400, "data": nil})
		return
	}
	utils.WriterJson(writer, map[string]interface{}{"err": "", "code": 200, "data": "success"})
}



func (User) SelectUser(writer http.ResponseWriter, request *http.Request) {
	users , err := model.Db.UserSelect()
	if err != nil{
		utils.WriterJson(writer, map[string]interface{}{"err": err.Error(), "code": 400, "data": nil})
		return
	}
	fmt.Println("users:",users)
	utils.WriterJson(writer, map[string]interface{}{"err": "", "code": 200, "data": users})
}