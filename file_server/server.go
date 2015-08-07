package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"unsafe"
)

const (
	ENV_FILE_SERVER_SECRET_KEY = "ENV_FILE_SERVER_SECRET_KEY"
	ENV_FILE_SERVER_TOKEN      = "ENV_FILE_SERVER_TOKEN"
)

var (
	// 文件根目录
	RootDir = "G:\\code\\go\\src\\github.com\\lab204\\mmDeploy\\file_server\\"

	// 用于文件上传认证
	Token     string = os.Getenv(ENV_FILE_SERVER_TOKEN)
	SecretKey string = os.Getenv(ENV_FILE_SERVER_SECRET_KEY)
)

type File struct {
	Name   string
	Path   string
	Force  bool
	Rename bool
}

func init() {
	if SecretKey == "" || Token == "" {
		SecretKey = "1"
		Token = "2"
		// panic("ENV CONF ERROR")
	}
}
func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Get请求
func Get(w http.ResponseWriter, res *http.Request) {
	if strings.ToUpper(res.Method) != "GET" {
		http.Error(w, `{"error":"Method Error"}`, http.StatusMethodNotAllowed)
		return
	}

	path := RootDir + res.URL.Path
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}

// 上传文件
func Push(w http.ResponseWriter, res *http.Request) {
	if strings.ToUpper(res.Method) != "POST" {
		http.Error(w, `{"error":"Method Error"}`, http.StatusMethodNotAllowed)
		// return
	}

	err := res.ParseForm()
	if err != nil {
		http.Error(w, `{"error":"ParseFrom Failed"}`, http.StatusBadRequest)
		// return
	}
	res.Header.Add("Access-Control-Allow-Origin", "*")

	// 上传文件认证
	// token := res.Form.Get("token")
	// salt := res.Form.Get("salt")
	// sum := res.Form.Get("sum")
	// if !auth(SecretKey, token, salt, sum) {
	// 	http.Error(w, `{"error":"Auth Failed"}`, http.StatusNotAcceptable)
	// 	// return
	// }

	// 接收文件
	path := RootDir + res.Form.Get("path")
	force := res.Form.Get("force")
	rename := res.Form.Get("rename")
	fmt.Println(path)
	fmt.Println(force)
	fmt.Println(rename)
	f, fh, err := res.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"ParseFileForm  Failed"}`, http.StatusBadRequest)
		// return
	}
	fmt.Println(fh.Filename)
	finfo, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0644)
	} else if !finfo.IsDir() {

	}
	newf, err := os.Create(path + fh.Filename)
	if err != nil {
		http.Error(w, `{"error":"Create File Failed"}`, http.StatusBadRequest)
	}
	_, err = io.Copy(newf, f)
	if err != nil {
		http.Error(w, `{"error":"Upload File Failed"}`, http.StatusBadRequest)
	}

	w.Write([]byte("In Push Method world" + res.Method + " <br> Form" + res.Form.Get("hello")))
}

func auth(sec, tok, sal, sum string) bool {
	md5sum := func(b []byte) string {
		h := md5.New()
		h.Write(b)
		return hex.EncodeToString(h.Sum(nil))
	}

	ss := []string{sec, tok, sal}
	sort.Strings(ss)
	return sum == md5sum([]byte(ss[0]+ss[1]+ss[2]))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Get)
	mux.HandleFunc("/push/", Push)

	fmt.Println("Http is running at 80")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}
}
