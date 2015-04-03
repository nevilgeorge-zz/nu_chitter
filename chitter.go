// NU Chitter - Project 0 - EECS 345: Distributed Systems - Northwestern University
// Nevil George - nsg622

package main

import (
	"fmt"
	"net"
	"os"
)

// // client struct to keep track of channels and id of each client
// type Client struct {
// 	id int
// 	incoming chan string
// 	outgoing chan string

// }

// // functions to be called on Client struct
// func (client *Client) Read() {
// 	// forever loop that keeps reading from the channel
// 	for {

// 	}
// }

func main() {
	var port string = os.Args[1]
	var count int = 0
	
	// listening for incoming connections
	listen, err := net.Listen("tcp", "localhost:" + port)
	if err != nil {
		fmt.Println("Error on listening: ", err.Error())
		os.Exit(1)
	}

	// close listening connection when function ends
	defer listen.Close()
	fmt.Println("Listening on localhost:" + port)

	// infinite loop waiting for incoming connections
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error occurred in accepting a new connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("New connection!")
		// no error, handle the new connection
		count += 1
		fmt.Println(count)
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	// byte slice required because net.Conn uses byte slices instead of strings
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			// fmt.Println("Error reading from connection: ", err.Error())
			// close connection when error occurs
			conn.Close()
			break
		}
		// Send a response to the client
		conn.Write(buffer)
	}
}
