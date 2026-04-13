package server

import (
	"net"
	"time"
)

func handleHello(req *Request) *Response {
	time.Sleep(2 * time.Second)
	return &Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       []byte("Hello, world!"),
	}
}

func handleEcho(req *Request) *Response {
	body := req.Body
	if len(body) == 0 {
		body = []byte("(no body)")
	}
	return &Response{
		StatusCode: 200,
		StatusText: "OK",
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       body,
	}
}

func notFound(req *Request) *Response {
	return &Response{
		StatusCode: 404,
		StatusText: "Not Found",
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       []byte("404 Not Found"),
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		WriteResponse(conn, &Response{
			StatusCode: 400,
			StatusText: "Bad Request",
			Headers:    map[string]string{"content-type": "text/plain"},
			Body:       []byte(err.Error()),
		})
		return
	}

	resp := route(req)
	WriteResponse(conn, resp)
}
