package Request

import (
	"strings"
)

type Request struct {
	raw []byte
	valid bool
	rType string
	path  string
	headersMap map[string]string
}

func NewRequest(raw []byte) Request {
	r := Request{}
	r.raw = raw
	r.valid = true
	r.headersMap = make(map[string]string)

	r.parseRaw()
	return r
}

func (r *Request) Path() string {
	return r.path
}

func (r *Request) parseRaw() {
	s := string(r.raw)
	s = strings.Replace(s, "\r", "", -1)
	headers := strings.Split(s, "\n")

	// Get request type
	if strings.Contains(headers[0], "GET") {
		r.rType = "GET"
	} else if strings.Contains(headers[0], "POST") {
		r.rType = "POST"
	} else if strings.Contains(headers[0], "HEAD") {
		r.rType = "HEAD"
	} else {
		r.valid = false
		return
	}

	// Get path
	path := strings.Split(s, " ")[1]
	path = strings.Split(path, "?")[0]
	if (strings.Index(path, "/") == 0) { // path must start with '/'
		r.path = path
	}

	// Map all other headers
	for _, v := range headers[1:] {
		if v == "" {
			break
		} // Stop when headers are no longer found
		keyEnd := strings.Index(v, ":")
		valueStart := strings.Index(v, " ") + 1

		key := v[:keyEnd]
		value := v[valueStart:]

		r.headersMap[key] = value
	}

	// Todo: Add parsing of GET parameters

}

