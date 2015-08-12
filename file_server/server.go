package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
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
	FILE_SERVER_SECRET_KEY = "FILE_SERVER_SECRET_KEY"
	FILE_SERVER_TOKEN      = "FILE_SERVER_TOKEN"
)

var (
	// 文件根目录
	RootDir = os.Getenv("DATA_DIR")

	// 用于文件上传认证
	Token = []byte(os.Getenv(FILE_SERVER_TOKEN))
	// SecretKey string = os.Getenv(FILE_SERVER_SECRET_KEY)
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
func Index(w http.ResponseWriter, req *http.Request) {
	defer func() {
		req.Body.Close()
	}()
	req.Header.Add("Access-Control-Allow-Origin", "*")
	err := req.ParseForm()
	if err != nil {
		http.Error(w, `{"error":"ParseFrom Failed"}`, http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	fmt.Println("func Get:", req.Method)
	switch strings.ToUpper(req.Method) {
	case "POST":
		Push(w, req)
		return
	case "GET":
	default:
		http.Error(w, `{"error":"Method Error"}`, http.StatusMethodNotAllowed)
		return
	}

	path := RootDir + req.URL.Path
	if !filter(path) {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	fif, err := os.Stat(path)
	if err != nil || fif.IsDir() {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}

// 上传文件
func Push(w http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)

	// 上传文件认证
	// token := req.Form.Get("token")
	// salt := req.Form.Get("salt")
	// sum := req.Form.Get("sum")
	// if !auth(SecretKey, token, salt, sum) {
	// 	http.Error(w, `{"error":"Auth Failed"}`, http.StatusNotAcceptable)
	// 	// return
	// }

	// 接收文件
	for k, v := range req.PostForm {
		fmt.Printf("Range # %s: %#v\n", k, v)
	}
	dir := RootDir + req.URL.Path + "/"
	force := false
	if req.Form.Get("force") == "true" {
		force = true
	}

	f, fh, err := req.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"ParseFileForm  Failed"}`, http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	//
	_, err = io.Copy(buf, f)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hsha1 := req.Header.Get("X-CKEYER-SHA1")
	if HmacSha1(buf.Bytes(), Token) != hsha1 {
		http.Error(w, `{"error":"Auth failed"}`, http.StatusNotAcceptable)
		return
	}

	path := dir + fh.Filename
	if !filter(path) {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	fmt.Println(fh.Filename)

	finfo, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	} else if !finfo.IsDir() {
		if force {
			os.Remove(dir)
			os.MkdirAll(dir, 0644)
		} else {
			http.Error(w, `{"error":"Dir is a exists File"}`, http.StatusBadRequest)
			return
		}
	}

	if _, err = os.Stat(path); err == nil {
		if force {
			os.Remove(path)
		} else {
			http.Error(w, `{"error":"File Exists"}`, http.StatusBadRequest)
			return
		}
	}

	newf, err := os.Create(path)
	if err != nil {
		http.Error(w, `{"error":"Create File Failed"}`, http.StatusBadRequest)
		return
	}
	defer newf.Close()
	_, err = io.Copy(newf, buf)
	if err != nil {
		http.Error(w, `{"error":"Upload File Failed"}`, http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"success":"ok"}`))
}

func filter(url string) bool {
	err_seps := []string{"..", "~", ".go", "--"}
	for _, sep := range err_seps {
		if strings.Count(url, sep) > 0 {
			return false
		}
	}
	return true
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
	mux.HandleFunc("/", Index)

	fmt.Println("Http is running at 80")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}
}

func HmacSha1(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return fmt.Sprintf("%x", expectedMAC)
}
