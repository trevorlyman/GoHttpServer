package Request

import (
	//"../Response"
)

type Request struct {
	raw []byte
	path  string
}

func NewRequest(raw []byte) Request {
	r := Request{}
	r.raw = raw
	r.parseRaw()
	return r
}

func (r Request) parseRaw() {
	// Todo: Parse the raw byte input
	r.path = "/"
}

