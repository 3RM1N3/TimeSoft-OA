package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var globalName, globalID string

type auth struct {
	Username string `json:"username"`
	Pwd      string `json:"password"`
}

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	ID   string `json:"id"`
}

func login(addr, username, password string) (Resp, error) {
	//post请求提交json数据
	auths := auth{username, password}
	bytesJson, err := json.Marshal(auths)
	if err != nil {
		return Resp{}, err
	}
	resp, err := http.Post(addr+"/login", "application/json", bytes.NewBuffer(bytesJson))
	if err != nil {
		return Resp{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Resp{}, err
	}

	var result Resp
	json.Unmarshal(body, &result)

	return result, nil
}

func get() {
	//get请求
	//http.Get的参数必须是带http://协议头的完整url,不然请求结果为空
	resp, _ := http.Get("http://localhost:8080/login2?username=admin&password=123456")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(body))
	fmt.Printf("Get request result: %s\n", string(body))
}

func postWithJson() {
	//post请求提交json数据
	auths := auth{"admin", "123456"}
	ba, _ := json.Marshal(auths)
	resp, _ := http.Post("http://localhost:8080/login1", "application/json", bytes.NewBuffer([]byte(ba)))
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Post request with json result: %s\n", string(body))
}

func postWithUrlencoded() {
	//post请求提交application/x-www-form-urlencoded数据
	form := make(url.Values)
	form.Set("username", "admin")
	form.Add("password", "123456")
	resp, _ := http.PostForm("http://localhost:8080/login2", form)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Post request with application/x-www-form-urlencoded result: %s\n", string(body))
}
