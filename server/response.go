package server

import (
	"fmt"
	"net"
	"strconv"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func WriteResponse(conn net.Conn, resp *Response) error {
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

	// Mandatory blank line
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
