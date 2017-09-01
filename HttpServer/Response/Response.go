package Response

import (
	"../Request"
	"net"
	"os"
	"strings"
	"mime"
	"bytes"
	"strconv"
	"fmt"
)

const WWW_FOLDER = "www"

type Response struct {
	headers string
	systemPath string
	responseCode string
	mimeType string
	contentLength int64
}


func NewResponse(request Request.Request) Response {
	// Todo: Look at the Request in order to construct the response
	_ = request
	r := Response{}


	code := "200"
	_ = code
	// What code to use?
	// 200, 403, 404, 500

	// can the file be found?
	r.systemPath = getSystemPath(request.Path())


	file, err := os.Open(r.systemPath)


	defer file.Close()
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			code = "404 Not Found"
		}else if strings.Contains(err.Error(), "ermission") {
			code = "403 Forbidden"
		}else {
			code = "500 Server Error"
		}
	}else {
		// Get content info

		//ext := r.systemPath[strings.LastIndex(r.systemPath, "."):]
		ext := getExt(r.systemPath)
		r.mimeType = mime.TypeByExtension(ext)
		fi, _ := file.Stat()
		r.contentLength = fi.Size()
	}

	r.responseCode = "HTTP/1.1 " + code
	contentLengthText := "Content-Length: " + strconv.Itoa(int(r.contentLength))

	// Build full headers
	buf := bytes.NewBufferString("")
	buf.WriteString(r.responseCode + "\r\n")
	buf.WriteString(contentLengthText + "\r\n")
	buf.WriteString("content-type: " + r.mimeType + "\r\n")
	buf.WriteString("\r\n")

	r.headers = buf.String()

	return r
}

func getSystemPath(requestPath string) string{
	ext := getExt(requestPath)
	var path string
	// Add ".html" to any request without an extension
	if requestPath == "/" {
		path = "/index.html"
	} else if ext == "" {
		path = requestPath + ".html"
	}
	return WWW_FOLDER + path

}

func getExt(s string) string {
	index := strings.LastIndex(s, ".")
	var ext string
	if index != -1 {
		ext = s[index:]
	}
	// else ""
	return ext
}

func (r *Response) Send(conn net.Conn) {

	file, err := os.Open(r.systemPath)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	// Send headers on the connection, read the file into the connection
	conn.Write([]byte(r.headers))
	//index := 0
	buf := make([]byte, 1024)

	n, err := file.Read(buf)
	for n != 0 {
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		s := string(buf)
		l := len(s)
		_ = s
		_ = l
		conn.Write(buf)
		n, err = file.Read(buf)
	}

}
