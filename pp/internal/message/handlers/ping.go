package handlers

import (
	"log"
	"pp/internal/message"
	// "pp/internal/node" //不再需要node
)

// PingHandler handles ping messages.
type PingHandler struct {
	serverAddr      string //需要serverAddr
	sendMessageFunc func(addr string, msg message.Message) error
}

// NewPingHandler creates a new PingHandler instance.
func NewPingHandler(serverAddr string, sendMessageFunc func(addr string, msg message.Message) error) *PingHandler {
	return &PingHandler{serverAddr: serverAddr, sendMessageFunc: sendMessageFunc}
}

// Handle processes a ping message.
func (h *PingHandler) Handle(senderAddr string, msg message.Message) {
	log.Printf("Received ping from %s", senderAddr)

	// Respond with a pong
	pongMsg := message.Message{Type: "pong", Data: "pong", Sender: h.serverAddr}
	if err := h.sendMessageFunc(senderAddr, pongMsg); err != nil { //不再需要node
		log.Printf("Error sending pong to %s: %v", senderAddr, err)
	}
}
