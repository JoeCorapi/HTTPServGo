package server

type HandlerFunc func(*Request) *Response

func route(req *Request) *Response {
	routes := map[string]HandlerFunc{
		"GET /hello": handleHello,
		"POST /echo": handleEcho,
	}

	key := req.Method + " " + req.Path
	if handler, ok := routes[key]; ok {
		return handler(req)
	}

	return notFound(req)
}
