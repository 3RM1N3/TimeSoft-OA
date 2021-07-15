package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	ID   string `json:"id"`
}

type Auth struct {
	Username string `json:"username"`
	Pwd      string `json:"password"`
}

//post接口接收json数据
func postLogin(writer http.ResponseWriter, request *http.Request) {
	var auth Auth

	if err := json.NewDecoder(request.Body).Decode(&auth); err != nil {
		request.Body.Close()
		log.Fatal(err)
	}

	var result Resp
	result.ID = ""

	if auth.Username == "admin" && auth.Pwd == "admin" {
		result.Code = "200"
		result.Msg = "登录成功"
		result.ID = "q65wer5tyu34i345kjhgfd2f3hgf4567"
		log.Printf("%s登录成功\n", auth.Username)
	} else {
		result.Code = "401"
		result.Msg = "账户名或密码错误"
	}

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}
}

//接收x-www-form-urlencoded类型的post请求或者普通get请求
func login2(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username, uError := request.Form["username"]
	pwd, pError := request.Form["password"]

	var result Resp
	if !uError || !pError {
		result.Code = "401"
		result.Msg = "登录失败"
	} else if username[0] == "admin" && pwd[0] == "admin" {
		result.Code = "200"
		result.Msg = "登录成功"
	} else {
		result.Code = "203"
		result.Msg = "账户名或密码错误"
	}

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}
}
