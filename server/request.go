package server

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
