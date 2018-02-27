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

const (
	// Todo folder const should be moved to the HttpServer level
	WWW_FOLDER = "www"
	MESSAGE_404_PATH = "messages/404.html"
	MESSAGE_403_PATH = "messages/403.html"
	MESSAGE_500_PATH = "messages/500.html"
)

type Response struct {
	request Request.Request
	headers string
	systemPath string
	responseCode string
	mimeType string
	contentLength int64
}


func NewResponse(request Request.Request) Response {
	r := Response{}
	r.request = request

	// Assume the request is 200 OK unless proven wrong
	code := "200 OK"


	// Can the file be found?
	r.systemPath = getSystemPath(request.Path())
	file, err := os.Open(r.systemPath)

	if err != nil {

		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			code = "404 Not Found"
			r.systemPath = MESSAGE_404_PATH
		}else if strings.Contains(strings.ToLower(err.Error()), "permission") {
			code = "403 Forbidden"
			r.systemPath = MESSAGE_403_PATH
		}else {
			code = "500 Server Error"
			r.systemPath = MESSAGE_500_PATH
		}

		file.Close()

		// If the html error file is not found
		file, err = os.Open(r.systemPath)
		if err != nil {
			fmt.Println("Error file not found: " + r.systemPath)

			// Todo: This is a bit harsh to close the program after not finding an error file,
			// refactor so that NewResponse returns (Response, error).
			os.Exit(1)
		}
	}


	// Get content info
	ext := getExt(r.systemPath)
	r.mimeType = mime.TypeByExtension(ext)
	fi, _ := file.Stat()
	r.contentLength = fi.Size()


	r.responseCode = "HTTP/1.0 " + code
	contentLengthText := "Content-Length: " + strconv.Itoa(int(r.contentLength))

	// Build full headers
	endl := "\r\n"

	buf := bytes.NewBufferString("")
	buf.WriteString(r.responseCode + endl)
	buf.WriteString(contentLengthText + endl)
	buf.WriteString("content-type: " + r.mimeType + endl)
	buf.WriteString(endl)

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
	} else {
		path = requestPath
	}
	return WWW_FOLDER + path

}

func getExt(s string) string {
	index := strings.LastIndex(s, ".")
	var ext string
	if index != -1 {
		ext = s[index:]
	}

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

	buf := make([]byte, 1024)

	n, err := file.Read(buf)
	for n != 0 {
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		 _, writeErr := conn.Write(buf)

		 if writeErr != nil {
			fmt.Print(writeErr.Error())
			break
		 }
		n, err = file.Read(buf)
	}

}
