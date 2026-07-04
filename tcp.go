package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const PORT = "127.0.0.1:6379"

func RunServer() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error setting up tcp socket: %v\n", err)
	}

	defer listener.Close()

	// single threaded channel that all connection commands are passed to
	commandCh := make(chan Command)
	go RunStore(commandCh)

	fmt.Println("Server is running on port:" + PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %v\n", err)
			continue
		}
		
		go handleConnection(conn, commandCh)
	}
}

func handleConnection(conn net.Conn, commands chan Command) {
	log.Printf("Accepted new connection from client")
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// the channel unique to each connection where replies are sent
	replyCh := make(chan Result, 1)

	for {
		cmd, err := ParseResp(reader, replyCh)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error parsing connection content: %v\n", err)
			WriteResp(conn, Result{err: err})
			continue
		}
		commands <- cmd

		result := <-replyCh
		WriteResp(conn, result)
	}
}
