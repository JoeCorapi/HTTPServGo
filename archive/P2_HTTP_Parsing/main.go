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

func parseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	// Requests come in 3 "Chunks"
	// Request MetaData, Headers, Body

	// --- Parse the request line ---
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

	// --- Parse headers ---
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		// Blank line = end of headers
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

	// --- Parse body if Content-Length present ---
	if lengthStr, ok := req.Headers["content-length"]; ok {
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid content-length: %w", err)
		}

		// Important: ReadFull handles the loop to read len(buf) into our buffer
		body := make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}

		req.Body = body
	}

	return req, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	fmt.Printf("Method:  %s\n", req.Method)
	fmt.Printf("Path:    %s\n", req.Path)
	fmt.Printf("Version: %s\n", req.Version)
	fmt.Println("Headers:")
	for k, v := range req.Headers {
		fmt.Printf("  %s: %s\n", k, v)
	}
	fmt.Printf("Body:    %q\n", req.Body)
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
		// Pull one connection off the OS queue
		// Blocks here until a client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		fmt.Println("Connection from:", conn.RemoteAddr())
		handleConnection(conn)
	}
}
