// Package challenge8 contains the solution for Challenge 8: Chat Server with Channels.
package challenge8

import (
	"errors"
	"sync"
	"fmt"
	// Add any other necessary imports
)

// Client represents a connected chat client
type Client struct {
	// TODO: Implement this struct
	// Hint: username, message channel, mutex, disconnected flag
	username     string
	messageChan  chan string
	mu           sync.Mutex
	disconnected bool
}

// Send sends a message to the client
func (c *Client) Send(message string) {
	// TODO: Implement this method
	// Hint: thread-safe, non-blocking send
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.disconnected {
		return
	}
	select {
	case c.messageChan <- message:
		// sent
	default:
		// if buffer full, drop message silently or handle accordingly
	}
}

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string {
	// TODO: Implement this method
	// Hint: read from channel, handle closed channel
	msg, ok := <-c.messageChan
	if !ok {
		return ""
	}
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	// TODO: Implement this struct
	// Hint: clients map, mutex
	clients map[string]*Client
	mu      sync.Mutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	// TODO: Implement this function
	return &ChatServer{
		clients: make(map[string]*Client),
	}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	// TODO: Implement this method
	// Hint: check username, create client, add to map
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	}

	client := &Client{
		username:    username,
		messageChan: make(chan string, 10), // buffered channel
		disconnected: false,
	}

	s.clients[username] = client
	s.Broadcast(nil, fmt.Sprintf("%s has joined the chat", username))
	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	// TODO: Implement this method
	// Hint: remove from map, close channels
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.disconnected {
		return
	}
	client.mu.Lock()
	client.disconnected = true
	client.mu.Unlock()

	delete(s.clients, client.username)
	close(client.messageChan)
	s.Broadcast(nil, fmt.Sprintf("%s has left the chat", client.username))
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	// TODO: Implement this method
	// Hint: format message, send to all clients
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	from := "Server"
// 	if sender != nil {
// 		from = sender.username
// 	}
// 	formatted := fmt.Sprintf("%s: %s", from, message)

// 	for _, client := range s.clients {
// 		if client != sender {
// 			client.Send(formatted)
// 		}
// 	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	// TODO: Implement this method
	// Hint: find recipient, check errors, send message
	s.mu.Lock()
	target, exists := s.clients[recipient]
	s.mu.Unlock()

	if !exists {
		return ErrRecipientNotFound
	}
	if sender.disconnected {
	    return ErrClientDisconnected
	}
	if target.disconnected {
	    return ErrClientDisconnected
	}

	if sender != nil {
		target.Send(fmt.Sprintf("[Private] %s: %s", sender.username, message))
	} else {
		target.Send(fmt.Sprintf("[Private] Server: %s", message))
	}
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
	// Add more error types as needed
)
