package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("error connecting to server: %v \n", err)
		return
	}
	defer conn.Close()
	fmt.Println("connected to server")

	for {
		fmt.Print("enter message: ")
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading inputted message: %v \n", err)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("error sending message: %v \n", err)
			return
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("error reading response: %v \n", err)
			return
		}
		fmt.Printf("reading response from server: %v \n", resp)
	}
}
