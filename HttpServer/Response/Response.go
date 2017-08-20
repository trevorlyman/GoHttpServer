package Response

import (
	"../Request"
)

type Response struct {
	Raw []byte
}


func NewResponse(request Request.Request) Response {
	// Todo: Look at the Request in order to construct the response
	_ = request
	r := Response{}
	r.Raw = []byte(`HTTP/1.1 200 OK
	Content-Length: 147

	<html>
	<h1>Hello-World!</h1>
	<h3>This is a go webserver, I hope you like it!</h3>
	<h4>If you don't like it then, that's too bad :(</h4>
</html>`)
	return r
}
