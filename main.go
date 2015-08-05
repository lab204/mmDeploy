package main

import (
	"fmt"
	"net/http"
)

//
func Index(w http.ResponseWriter, res *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("hello world"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	fmt.Println("Http is running")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
