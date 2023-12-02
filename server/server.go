package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// ChatServer represents the chat server with nickname management and broadcasting.
type ChatServer struct {
	clients map[string]net.Conn
	mutex   sync.Mutex
}

func (s *ChatServer) handleClientCommand(conn net.Conn, data string) {
	// Handle the command logic here
	fields := strings.Fields(data)
	if len(fields) == 0 {
		// Handle the case where no command is provided
		fmt.Fprintf(conn, "Error: No command provided\n")
		return
	}
	command := fields[0]
	args := fields[1:]

	switch command {
	case "/LIST":
		go s.handleListCommand(conn)
	case "/NICK":
		go s.handleNickCommand(conn, args)
	case "/BC":
		go s.handleBcCommand(conn, args)
	case "/MSG":
		go s.handleMsgCommand(conn, args)
	default:
		go s.handleUnknownCommand(conn)
	}
}

func (s *ChatServer) handleListCommand(conn net.Conn) {

	resultChan := make(chan string)
	go s.createNameList(resultChan)
	namesList := <-resultChan

	fmt.Fprintf(conn, "List of NickNames: %s\n", namesList)
}

func (s *ChatServer) handleNickCommand(conn net.Conn, args []string) {

	if len(args) != 1 {
		// Invalid usage of /NICK command
		fmt.Fprintf(conn, "Error: Invalid usage of /NICK command. Usage: /NICK <nickname>\n")
		return
	}

	nickname := args[0]

	if !s.validateNickname(nickname) {
		fmt.Fprintf(conn, "Error: Invalid nickname format. \n")
		return
	}
	go s.registerNickname(conn, nickname)

}

func (s *ChatServer) handleBcCommand(conn net.Conn, args []string) {

	nickname, found := s.getNicknameByConn(conn)
	if !found {
		fmt.Fprintf(conn, "Error: You must register a nickname before broadcasting.\n")
		return
	}
	message := strings.Join(args, " ")
	s.broadcastMessage(nickname, message)
}

func (s *ChatServer) handleMsgCommand(conn net.Conn, args []string) {

	if len(args) < 2 {
		fmt.Fprintf(conn, "Error: Invalid usage of /MSG command. Usage: /MSG <nickname> <message>\n")
		return
	}

	receiverNickname := args[0]
	message := args[1:]

	nickname, found := s.getNicknameByConn(conn)
	if !found {
		fmt.Fprintf(conn, "Error: You must register a nickname before broadcasting.\n")
		return
	}

	receiverConn, err := s.findReceiverConnection(conn, receiverNickname)
	if err != nil {
		return
	}

	fmt.Fprintf(receiverConn, "[%s]: %s\n", nickname, strings.Join(message, " "))
}

func (s *ChatServer) handleUnknownCommand(conn net.Conn) {
	// Implement handling for unknown command
	fmt.Fprintf(conn, "Unknown command\n")
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]net.Conn),
		mutex:   sync.Mutex{},
	}
}

func main() {
	server := NewChatServer()

	ln, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		// server := &ChatServer{} // Create an instance of ChatServer
		go server.handleClient(conn)
	}

}

func (s *ChatServer) handleClient(conn net.Conn) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	// fmt.Fprintf(conn, "%s\n", sc.Text())

	for sc.Scan() {
		text := sc.Text()
		go s.handleClientCommand(conn, text)
	}

	go s.handleClientDisconnected(conn)

}

func (s *ChatServer) handleClientDisconnected(conn net.Conn) {
	// Perform cleanup or handle disconnection logic here
	nickname, _ := s.getNicknameByConn(conn)
	delete(s.clients, nickname)

}
