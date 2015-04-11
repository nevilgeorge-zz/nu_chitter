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
	id       int
	incoming chan string
	outgoing chan string
	conn     net.Conn
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

// message struct to keep track of sender id, receiver id and text
type Message struct {
	sender_id int
	receiver_id int
	text string
}

// array to keep track of clients currently connected to the server
// var clients = make([]*Client, 0)
var clientChannel chan *Client
var countChannel chan int
var count int = 0
var broadcastChannel chan *Message
var messageChannel chan *Message

func main() {
	var port string = os.Args[1]

	// listening for incoming connections
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error on listening: ", err.Error())
		os.Exit(1)
	}

	clientChannel = make(chan *Client)
	countChannel = make(chan int)
	broadcastChannel = make(chan *Message)
	messageChannel = make(chan *Message)
	go handleClients()

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
		clientChannel <- newClient
		go handleRequest(conn, *newClient)
	}
}

func handleClients() {
	var clients = make([]*Client, 0)
	for {
		select {
		case newCli := <- clientChannel:
			clients = append(clients, newCli)
			// countChannel <- len(clients)

		case newMsg := <- broadcastChannel:
			broadcastMsg(clients, newMsg)

		case newPM := <- messageChannel:
			sendMessage(clients, newPM)
		}
	}
}

func handleRequest(conn net.Conn, cli Client) {
	// create separate goroutines to read from and write to the client
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
			} else if isPersonalMessage(msgIn) {
				// get id of receiver, grab message text and then send the message to the recipient
				id, err := strconv.Atoi(string(msgIn[0]))
				if err != nil {
					continue
				}
				text := grabTextAfterColon(msgIn)
				newMessage := new(Message)
				newMessage.text = text
				newMessage.sender_id = cli.id
				newMessage.receiver_id = id
				messageChannel <- newMessage
			} else if msgIn[0:3] == "all" {
				// grab message text and broadcast to all connected clients
				text := grabTextAfterColon(msgIn)
				newMessage := new(Message)
				newMessage.text = text
				newMessage.sender_id = cli.id
				broadcastChannel <- newMessage
			} else {
				// broadcast message by default
				newMessage := new(Message)
				newMessage.text = msgIn
				newMessage.sender_id = cli.id
				broadcastChannel <- newMessage
			}
		}
	}
}

// function that checks whether the first entry is a number of this
func isPersonalMessage(input string) bool {
	text := strings.TrimSpace(input)
	_, err := strconv.Atoi(string(text[0]))
	index := strings.Index(input, ":")
	if err == nil && index > -1 {
		return true
	} else {
		return false
	}
}

// function that sends a message to all connected clients
// params: msg - string
func broadcastMsg(clients []*Client, msg *Message) {
	for x := 0; x < len(clients); x++ {
		text := strconv.Itoa(msg.sender_id) + " : " + msg.text
		clients[x].outgoing <- text
	}
}

// function that parses the string and grabs the message being sent
// params: input - string, return: text - string
func grabTextAfterColon(input string) (text string) {
	stringArr := strings.Split(input, ":")
	text = stringArr[1]
	return strings.TrimSpace(text)
}

// function to send a personal message to a client
func sendMessage(clients []*Client, msg *Message) {
	chat := strconv.Itoa(msg.sender_id) + " : " + msg.text
	clients[msg.receiver_id - 1].outgoing <- chat

}

