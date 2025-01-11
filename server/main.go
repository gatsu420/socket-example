package main

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
)

var clientCounter uint32

func main() {
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
		clientID := atomic.AddUint32(&clientCounter, 1)
		fmt.Printf("client %v connected \n", clientID)

		go handleConnection(conn, clientID)
	}
}

func handleConnection(conn net.Conn, clientID uint32) {
	defer conn.Close()

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
			fmt.Printf("error sending response message to client%v: %v \n", clientID, err)
			return
		}
	}
}
