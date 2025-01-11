package main

import (
	"bufio"
	"fmt"
	"net"
)

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
		fmt.Println("client connected")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("client disconnected")
			return
		}
		fmt.Printf("received from client: %v \n", message)

		resp := fmt.Sprintf("received message %v", message)
		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Printf("error sending response message: %v \n", err)
			return
		}
	}
}
