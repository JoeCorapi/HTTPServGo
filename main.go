package main

import (
	"fmt"
	"httpservgo/server"
	"net"
)

func main() {

	// Initializer server worker pool
	pool := server.NewWorkerPool(3, 6) // 3 workers, queue holds 6

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to bind:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :8080 (workers: 3, queue: 6)")

	// Sempahore to limit number active threads
	// Remove semaphore for concrete implementation - not go channels
	// sem := make(chan struct{}, 10) // max 10 concurrent requests

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		if !pool.Submit(conn) {
			fmt.Println("Pool saturated — rejecting connection")
			server.WriteResponse(conn, &server.Response{
				StatusCode: 503,
				StatusText: "Service Unavailable",
				Headers:    map[string]string{"content-type": "text/plain"},
				Body:       []byte("server is at capacity"),
			})
			conn.Close()
		}
	}
}
