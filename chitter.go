// NU Chitter - Project 0 - EECS 345: Distributed Systems - Northwestern University
// Nevil George - nsg622

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// client struct to keep track of channels and id of each client
type Client struct {
	id int
	incoming chan string
	outgoing chan string
	conn net.Conn
}

// methods for the Client struct
func (client *Client) Read() {
	conn := client.conn
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			break
		}
		input := string(buffer)
		client.incoming <- input
	}
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.conn.Write([]byte(data))
	}
}

func (client *Client) PrintId() {
	t := strconv.Itoa(client.id)
	t += "\n"
	client.outgoing <- t
}

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
		// no error, handle the new connection
		count += 1
		fmt.Println("New connection! Id: " + strconv.Itoa(count))
		// create a new Client instance
		newClient := new(Client)
		newClient.id = count
		newClient.incoming = make(chan string)
		newClient.outgoing = make(chan string)
		newClient.conn = conn
		clients = append(clients, newClient)
		go handleRequest(conn, *newClient)
	}

}

func handleRequest(conn net.Conn, cli Client) {
	// create separate go routines to read from and write to the client
	go cli.Read()
	go cli.Write()
	// run the select statement in the current go routine
	for {
		select {
			// when a incoming message is present, we need to parse it and then handle it correctly
		case msgIn := <-cli.incoming:
			if msgIn[0:6] == "whoami" {
				// print id
				cli.PrintId()
			} else if msgIn[0:3] == "all" {
				// grab message text and broadcast to all connected clients
				text := grabTextAfterColon(msgIn)
				broadcastMsg(text)
			} else {
				// get id of receiver, grab message text and then send the message to the recipient
				id, err := strconv.Atoi(string(msgIn[0]))
				if err != nil || id > len(clients) {
					continue
				}
				text := grabTextAfterColon(msgIn)
				sendToRecipient(cli.id, id, text)
			}
		}
	}
}

// function that sends a message to all connected clients
// params: msg - string
func broadcastMsg(msg string) {
	for x := 0; x < len(clients); x++ {
		clients[x].outgoing <- msg
	}
}

// function that parses the string and grabs the message being sent
// params: input - string, return: text - string
func grabTextAfterColon(input string) (text string) {
	index := strings.Index(input, ":")
	if input[index + 1] == ' ' {
		text = input[index + 2:]
	} else {
		text = input[index + 1:]
	}
	return text
}

// function that handles sending a private message 
// params: sender id  - int, receiver id - int, text - string
func sendToRecipient(sender int, receiver int, text string) {
	message := strconv.Itoa(sender) + " : " + text
	clients[receiver - 1].outgoing <- message
}