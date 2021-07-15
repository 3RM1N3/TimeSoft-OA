package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/login", postLogin)
	//http.HandleFunc("/login2", login2)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
