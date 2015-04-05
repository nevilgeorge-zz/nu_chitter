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
	message chan string
	conn net.Conn
}

func (client *Client) Write(msg string) {
	client.conn.Write([]byte(msg))
}

func (client *Client) PrintId() {
	t := strconv.Itoa(client.id)
	t += "\n"
	client.conn.Write([]byte(t))
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
		input := string(buffer)
		if input[0:6] == "whoami" {
			cli.PrintId()
		} else if input[0:3] == "all" {
			text := grabTextAfterColon(input)
			broadcastMsg(text)
		} else {
			id, err := strconv.Atoi(string(input[0]))
			if err != nil {
				continue
			}
			text := grabTextAfterColon(input)
			sendToRecipient(id, text)
		}
		
		// broadcast <- string(buffer)
		select {
		case personalMsg := <-cli.message:
			conn.Write([]byte(personalMsg))
		case broadcastMsg := <-broadcast:
			conn.Write([]byte(broadcastMsg))
		default:
			// we need to have a default otherwise the select statement waits
			// just continue iterating through the for loop
			continue
		}
	}
}

func broadcastMsg(msg string) {
	for x := 0; x < len(clients); x++ {
		clients[x].Write(msg)
	}
}

func grabTextAfterColon(input string) (text string) {
	index := strings.Index(input, ":")
	if input[index + 1] == ' ' {
		text = input[index + 2:]
	} else {
		text = input[index + 1:]
	}
	return text
}

func sendToRecipient(id int, text string) {
	curr := clients[id - 1]
	curr.Write(text)
}