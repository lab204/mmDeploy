package main

import (
	"fmt"
	"net/http"
)

//
func Push(w http.ResponseWriter, res *http.Request) {
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("hello world"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", Push)
	
	fmt.Println("Http is running at 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
