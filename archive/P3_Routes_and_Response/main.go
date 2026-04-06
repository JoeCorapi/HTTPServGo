package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func parseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request line: %w", err)
	}

	requestLine = strings.TrimRight(requestLine, "\r\n")
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: %q", requestLine)
	}

	req := &Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
		Headers: make(map[string]string),
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			return nil, fmt.Errorf("malformed header: %q", line)
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])
		req.Headers[strings.ToLower(key)] = value
	}

	if lengthStr, ok := req.Headers["content-length"]; ok {
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid content-length: %w", err)
		}

		body := make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}

		req.Body = body
	}

	return req, nil
}

func writeResponse(conn net.Conn, resp *Response) error {
	// Status line
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.StatusCode, resp.StatusText)
	_, err := conn.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	// Always set Content-Length
	resp.Headers["content-length"] = strconv.Itoa(len(resp.Body))

	// Headers
	for k, v := range resp.Headers {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	// Blank line
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	// Body
	if len(resp.Body) > 0 {
		_, err = conn.Write(resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

// --- Handlers ---

func handleHello(req *Request) *Response {
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

// --- Router ---

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

// --- Connection handler ---

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		writeResponse(conn, &Response{
			StatusCode: 400,
			StatusText: "Bad Request",
			Headers:    map[string]string{"content-type": "text/plain"},
			Body:       []byte(err.Error()),
		})
		return
	}

	resp := route(req)
	writeResponse(conn, resp)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to bind:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		fmt.Println("Connection from:", conn.RemoteAddr())
		handleConnection(conn)
	}
}
