// NU Chitter - Project 0 - EECS 345: Distributed Systems - Northwestern University
// Nevil George - nsg622

package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

// client struct to keep track of channels and id of each client
type Client struct {
	id int
	message chan string
	reader *bufio.Reader
	writer *bufio.Writer
	conn net.Conn
}

// functions to be called on Client struct
// func (client *Client) Read() {
// 	// forever loop that keeps reading from the channel
// 	for {
// 		line, _ := client.reader.ReadString('\n')
// 		// fmt.Println(line)
// 		client.incoming <- line
// 	}
// }

func (client *Client) Write(msg string) {
	// for data := range client.incoming {
	// 	fmt.Println(data)
	// 	client.writer.WriteString(data)
	// 	client.writer.Flush()
	// }
	client.conn.Write([]byte(msg))
}

// func (client *Client) Listen() {
// 	go client.Read()
// 	go client.Write()
// }

// func NewClient(connection net.Conn) *Client {
// 	writer := bufio.NewWriter(connection)
// 	reader := bufio.NewReader(connection)

// 	client := &Client {
// 		incoming: make(chan string),
// 		outgoing: make(chan string),
// 		reader: reader,
// 		writer: writer,
// 	}

// 	client.Listen()

// 	return client;
// }

// array to keep track of clients currently connected to the server
var clients = make([]*Client, 0)
var broadcast = make(chan string)

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
		newClient := new(Client)
		newClient.id = count
		newClient.message = make(chan string)
		newClient.conn = conn
		clients = append(clients, newClient)
		fmt.Println(len(clients))
		go handleRequest(conn, *newClient)
	}

}

func handleRequest(conn net.Conn, cli Client) {
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			break
		}
		// Broadcast(string(buffer))
		broadcast <- string(buffer)
		select {
		case personalMsg := <-cli.message:
			conn.Write([]byte(personalMsg))
		case broadcastMsg := <-broadcast:
			conn.Write([]byte(broadcastMsg))
		default:
			continue
		}
	}
}

func Broadcast(msg string) {
	for x := 0; x < len(clients); x++ {
		clients[x].Write(msg)
	}
}

// func handleRequest(conn net.Conn) {
// 	for {
// 	// byte slice required because net.Conn uses byte slices instead of strings
// 		buffer := make([]byte, 1024)
// 		_, err := conn.Read(buffer)
// 		if err != nil {
// 			// fmt.Println("Error reading from connection: ", err.Error())
// 			// close connection when error occurs
// 			conn.Close()
// 			break
// 		}
// 		// Send a response to the client
// 		conn.Write(buffer)
// 	}
// }
