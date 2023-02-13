// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	clientChannel chan<- string // an outgoing message channel
	clientName    string        // holds the name of the client
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	i := 0                           // Integer to hold client count
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.clientChannel <- msg
			}

		case cli := <-entering:
			clients[cli] = true

			if i > 0 { // Outputs message to clients if there are more than one client

				cli.clientChannel <- "Welcome! The current set of clients are: "

			}

			for temp := range clients {

				if temp != cli { // Output clients if different from entering client

					cli.clientChannel <- temp.clientName
					i++

				} else if i == 0 { // Output message if only one client

					cli.clientChannel <- "Welcome! You are the only client currently here"
					i++

				} else { // If temp matches current client, don't list

				}

			}

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.clientChannel)
			i-- // Decrement client count upon leaving
		}
	}
}

func handleConn(conn net.Conn) {

	ch := make(chan string) // outgoing client messages
	var cli client          // Client that will be added into map

	go clientWriter(conn, ch)

	ch <- "Enter your name (Server ID Name): " // Asks user to enter a name

	enterName := bufio.NewScanner(conn) // Variable to scan this connection

	enterName.Scan()
	cli.clientName = enterName.Text() // Assigns scanned name as client name

	who := cli.clientName // Who is assigned with user made client name
	ch <- "You are " + who
	messages <- who + " has arrived"
	cli.clientChannel = ch // Client channel is assigned with newly made channel
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli // Client leaves server and map
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {

	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
