// NU Chitter - Project 0 - EECS 345: Distributed Systems - Northwestern University
// Nevil George - nsg622

package main

import (
	"fmt"
	"net"
	"os"
)

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
	// byte array required because net.Conn uses byte arrays instead of strings
	buffer := make([]byte, 1024)
	request, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
	}
	fmt.Println(request)
	// Send a response to the client
	conn.Write([]byte("Message received."))
	// Close connection
	defer conn.Close()
}
