package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	url   = "http://ckeyer-file-s.daoapp.io/"
	Token = "123"
)

func main() {
	Upload()
}

func Upload() (err error) {
	// Create buffer
	buf := new(bytes.Buffer) // caveat IMO dont use this for large files, \
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)
	// Create file field
	fw, err := w.CreateFormFile("file", "test.jpg") //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		fmt.Println("c")
		return err
	}
	fd, err := os.Open("F:\\Pic\\beijing\\DSC03917.JPG")
	if err != nil {
		fmt.Println("d")
		return err
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		fmt.Println("e")
		return err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST", url+"test?force=true", buf)
	if err != nil {
		fmt.Println("f")
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-CKEYER-SHA1", HmacSha1(buf.Bytes(), Token))

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		return err
	}
	io.Copy(os.Stderr, res.Body) // Replace this with Status.Code check
	fmt.Println("h")
	return err
}

func HmacSha1(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return fmt.Sprintf("%x", expectedMAC)
}
