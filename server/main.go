package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type server struct {
	clientUUIDs map[string]int

	clientCounter atomic.Uint32
	mu            sync.Mutex
}

func main() {
	srv := &server{
		clientUUIDs: map[string]int{},
	}

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

	clientUUID := uuid.NewString()

	s.mu.Lock()
	s.clientUUIDs[clientUUID] = 0
	fmt.Printf("client with uuid %v connected \n", clientUUID)
	fmt.Printf("there are %v clients: %#v \n", len(s.clientUUIDs), s.clientUUIDs)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clientUUIDs, clientUUID)
		fmt.Printf("client with uuid %v disconnected \n", clientUUID)
		fmt.Printf("there are %v clients: %#v \n", len(s.clientUUIDs), s.clientUUIDs)
		s.mu.Unlock()
	}()

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
