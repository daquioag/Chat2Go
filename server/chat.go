package main

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

func (s *ChatServer) validateNickname(name string) bool {
	// Define the regular expression pattern
	pattern := "^[a-zA-Z][a-zA-Z0-9_]{0,11}$"
	regex := regexp.MustCompile(pattern)
	// Check if the nickname matches the pattern and is not empty after trimming
	return regex.MatchString(name) && strings.TrimSpace(name) != ""
}

func (s *ChatServer) createNameList(resultChan chan<- string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var namesList []string
	for nickname := range s.clients {
		namesList = append(namesList, nickname)
	}
	resultChan <- strings.Join(namesList, ", ")
}

func (s *ChatServer) isNicknameTaken(nickname string, conn net.Conn) bool {

	if userConn, exists := s.clients[nickname]; exists {
		if conn == userConn {
			fmt.Fprintf(conn, "You already have the nickname: '%s'\n", nickname)
		} else {
			fmt.Fprintf(conn, "Error: The nickname '%s' is already taken\n", nickname)
		}
		return true
	}
	return false
}

func (s *ChatServer) registerNickname(conn net.Conn, nickname string) {
	if s.isNicknameTaken(nickname, conn) {
		return
	}

	oldNickname, _ := s.getNicknameByConn(conn)
	delete(s.clients, oldNickname)
	s.clients[nickname] = conn
	fmt.Fprintf(conn, "You now have the nickname: '%s'\n", nickname)

}

func (s *ChatServer) broadcastMessage(senderNickname string, message string) {

	for nickname, clientConn := range s.clients {
		if nickname != senderNickname {
			// Send the broadcast message to all clients except the sender
			fmt.Fprintf(clientConn, "[%s]: %s\n", senderNickname, message)
		}
	}
}

func (s *ChatServer) getNicknameByConn(conn net.Conn) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for nickname, clientConn := range s.clients {
		if clientConn == conn {
			return nickname, true
		}
	}
	return "", false
}

func (s *ChatServer) findReceiverConnection(conn net.Conn, receiverNickname string) (net.Conn, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if receiverConn, exists := s.clients[receiverNickname]; exists {
		return receiverConn, nil
	}

	fmt.Fprintf(conn, "The name: %s does not exist!\n", receiverNickname)
	return nil, errors.New("receiver not found")
}
