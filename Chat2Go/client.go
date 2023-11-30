package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	serverMessages := make(chan string)

	go func() {
		sc := bufio.NewScanner(conn)
		for sc.Scan() {
			serverMessages <- sc.Text()

		}
	}()

	// Goroutine for reading messages from the keyboard
	go func() {
		kbd := bufio.NewScanner(os.Stdin)
		for {
			fmt.Printf("> ")
			if !kbd.Scan() {
				break
			}
			fmt.Fprintf(conn, "%s\n", kbd.Text())
		}
	}()

	// Main goroutine for handling messages
	for msg := range serverMessages {
		// Handle messages from the server
		fmt.Printf("%s\n", msg)
		fmt.Printf("> ")
	}

}
