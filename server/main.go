package main

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
)

type server struct {
	clientCounter atomic.Uint32
}

func main() {
	srv := &server{}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("error starting server: %v \n", err)
		return
	}
	defer listener.Close()
	fmt.Println("server is listening port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v \n", err)
			continue
		}

		go srv.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientID := s.clientCounter.Add(1)
	fmt.Printf("client %v connected \n", clientID)

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("client %v disconnected \n", clientID)
			return
		}
		fmt.Printf("received from client %v: %v \n", clientID, message)

		resp := fmt.Sprintf("client %v received message %v", clientID, message)
		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Printf("error sending response message to client: %v \n", err)
			return
		}
	}
}
