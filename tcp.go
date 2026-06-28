package redis

import (
	"fmt"
	"log"
	"net"
)

const PORT = "6379"

func RunServer() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error setting up tcp socket: %v\n", err)
	}

	defer listener.Close()

	fmt.Println("Server is running on port:" + PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %v\n", err)
			continue
		}
		
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

}
