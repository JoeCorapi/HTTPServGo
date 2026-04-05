package main

import (
	"fmt"
	"net"
)

func main() {
	// Ask the OS for a TCP socket, bound to port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to bind:", err)
		return
	}
	// Defers until enclosing method returns, similiar to catch
	defer listener.Close()

	fmt.Println("Listening on :8080")

	for {
		// Pull one connection off the OS queue
		// Blocks here until a client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err)
			continue
		}

		fmt.Println("Connection from:", conn.RemoteAddr())

		// Read raw bytes from the connection
		// The OS tracks data read from the kernel buffer ->
		// Calling conn.Read advances it's internal cursor so calling continually results in reads in sequence
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read:", err)
			conn.Close()
			continue
		}

		// Print exactly what arrived — raw bytes as a string
		fmt.Println("--- Raw request ---")
		fmt.Println(string(buf[:n]))
		fmt.Println("--- End ---")

		conn.Close()
	}
}
